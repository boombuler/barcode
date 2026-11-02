//go:build go1.26

package pdf417

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/boombuler/barcode"
)

// Scan tests render each fixture to a PNG and decode it with an
// external scanner to increase confidence that the encoded data is
// read exactly by the scanner.
//
// The BARCODE_SCANNERS environment variable configures which scanners
// are used by tests.
//
//   - auto (default): uses all scanners; skips uninstalled or old versions
//   - all: uses all known scanners; fails on uninstalled or old versions
//   - none: use no scanners; skips all scanner tests
//
// Alternatively, a comma-separated list of scanners can be provided
// from the following list:
//
//   - zxing: scans with ZXingReader from the zxing-cpp project.
//
// Skipped scanners report the reason when running tests with verbose output.
const scannersEnv = "BARCODE_SCANNERS"

// knownScanners lists the names accepted in BARCODE_SCANNERS.
var knownScanners = []string{"zxing"}

type scanTest struct {
	name  string
	data  string
	level byte
	// skip is set with a description about why a scan is expected to fail.
	skip string
}

var scanTests = []scanTest{
	// Examples from the specifications.
	// The standard leaves the layout to the encoder, so they are not checked here.
	{name: "spec/ISO-Figure1", data: "PDF417 Symbology Standard", level: 1}, // ISO/IEC 15438:2015 §5.2.1
	{name: "spec/ISO-5.4.2.2", data: "PDF417", level: 2},                    // text compaction example
	{name: "spec/ISO-5.4.4.2", data: "00246812345678", level: 2},            // numeric compaction example
	{name: "spec/AIM-Figure1", data: "AIM USA  PA", level: 1},               // AIM USA PDF417 §2.1
	// AIM USA PDF417 Appendix B. The library's symbol matches the figure
	// module for module, so it is a golden test candidate.
	{name: "spec/AIM-FigureB1", data: "AIM USA\r\n634 Alpha Drive\r\nPittsburgh, PA 15238", level: 2},

	// Text compaction sub-mode transitions.
	{name: "text/upper", data: "HELLO WORLD", level: 2}, // initial sub-mode, no latches
	{name: "text/upper-to-lower-latch", data: "hello world", level: 2},
	{name: "text/lower-to-upper-shift", data: "iPhone", level: 2},
	{name: "text/upper-to-mixed-latch", data: "AB-12", level: 2},
	{name: "text/mixed-to-upper-latch", data: "A1B", level: 2},
	{name: "text/mixed-to-lower-latch", data: "a1b", level: 2},
	{name: "text/upper-punct-shift", data: "A;B", level: 2},
	{name: "text/lower-punct-shift", data: "a;b", level: 2},
	{name: "text/mixed-punct-shift", data: "1;A", level: 2}, // shift, because the next character isn't punctuation
	{name: "text/mixed-to-punct-latch", data: "1((", level: 2},
	// Leaving punctuation. 39d842a, which is not in any release, dropped the
	// continue after the upper latch, losing the next character.
	{name: "text/punct-latch-leave", data: "1((A", level: 2},
	{name: "text/json", data: `{"a":[1,2],"b":3}`, level: 2},
	// Leaves punctuation through the odd-length pad, where 29 is the upper
	// latch, then a byte shift. Present since dbbaa13 (v1.0.0).
	{name: "text/punct-leave-padding", data: "0;;;;\x80abcdefghij", level: 2},

	// Numeric compaction.
	{name: "numeric/between-text", data: "ABC1234567890123DEF", level: 2},        // 13 digits: the shortest run that latches
	{name: "numeric/50-digits", data: strings.Repeat("1234567890", 5), level: 2}, // more than one 44-digit chunk

	// Byte compaction.
	{name: "byte/sixpack", data: "\x00\x01\x02\x80\xfe\xff", level: 2},    // latch 924
	{name: "byte/padded", data: "\x00\x01\x02\x80\xfe\xff\x7f", level: 2}, // latch 901
	{name: "byte/shift", data: "ABCDE\x01FGHIJ", level: 2},                // shift 913
	{name: "byte/utf8", data: "héllo wörld", level: 2},

	// The lowest and highest security levels.
	{name: "level/0", data: "Security level", level: 0},
	{name: "level/8", data: "Security level", level: 8},

	// The smallest symbol. ISO/IEC 15438:2015 §5.2.1 requires at least 3 rows.
	{name: "dims/min-rows", data: "A", level: 0,
		skip: "minRows = 2 gives 2 rows, which ZXingReader can't find"},
}

func TestScanZXing(t *testing.T) {
	reader := lookupZXing(t)
	artifacts := t.ArtifactDir()

	for _, tt := range scanTests {
		t.Run(tt.name, func(t *testing.T) {
			bc, err := Encode(tt.data, tt.level)
			if err != nil {
				t.Fatalf("Encode(%q, %d): %v", tt.data, tt.level, err)
			}
			if content := bc.Content(); content != tt.data {
				t.Errorf("Encode(%q, %d).Content() = %q, want %q", tt.data, tt.level, content, tt.data)
			}

			pngData, err := encodeScanPNG(bc)
			if err != nil {
				t.Fatalf("encodeScanPNG: %v", err)
			}
			pngPath := filepath.Join(artifacts, strings.ReplaceAll(tt.name, "/", "_")+".png")
			if err := os.WriteFile(pngPath, pngData, 0o644); err != nil {
				t.Fatal(err)
			}
			defer func() {
				if t.Failed() {
					t.Logf("Encode(%q, %d)\nscan image: %s (kept by go test -artifacts)",
						tt.data, tt.level, pngPath)
				}
			}()

			got, err := scanZXing(t, reader, pngPath)
			if err == nil && !bytes.Equal(got, []byte(tt.data)) {
				err = fmt.Errorf("decoded %q (% x), want %q (% x)", got, got, tt.data, tt.data)
			}
			switch {
			case err != nil && tt.skip != "":
				t.Skipf("ZXingReader: %v (known bug: %s)", err, tt.skip)
			case err != nil:
				t.Errorf("ZXingReader: %v", err)
			case tt.skip != "":
				t.Errorf("ZXingReader: decoded %q; known bug %q no longer reproduces, so remove skip", got, tt.skip)
			}
		})
	}
}

func encodeScanPNG(bc barcode.Barcode) ([]byte, error) {
	// Minimum width of a module (X) is up to application to
	// decide based on equipment for producing and reading
	// symbols. In our case, the scanner is the main constraint.
	const moduleWidth = 2 // px

	b := bc.Bounds()
	scaled, err := barcode.Scale(bc,
		b.Dx()*moduleWidth,
		b.Dy()*moduleWidth)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, scaled); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// scanZXing returns the content of the single PDF417 symbol in the path.
func scanZXing(t *testing.T, reader, path string) ([]byte, error) {
	t.Helper()

	// ZXingReader prints JSON for a symbol it finds, with Error set when the
	// symbol fails to decode. Otherwise stdout is empty, and the exit status and
	// stderr say why.
	cmd := exec.Command(reader, "-formats", "pdf417", "-single", "-errors", "-json", path)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	stdout, err := cmd.Output()
	if s := strings.TrimSpace(stderr.String()); s != "" {
		t.Logf("ZXingReader: stderr: %s", s)
	}
	if len(stdout) == 0 && err != nil {
		return nil, err
	}
	if len(stdout) == 0 {
		return nil, errors.New("no symbol found")
	}

	var res zxingResult
	if err := json.Unmarshal(stdout, &res); err != nil {
		return nil, fmt.Errorf("parsing output: %v\n%s", err, stdout)
	}
	if res.Error != "" {
		return nil, errors.New(res.Error)
	}
	if res.Format != "PDF417" {
		return nil, fmt.Errorf("format = %q, want PDF417", res.Format)
	}
	return res.Bytes, nil
}

// zxingResult holds the fields of ZXingReader -json output.
type zxingResult struct {
	Format string `json:"Format"`
	// Bytes is the exact for byte decoded content.
	Bytes hexBytes `json:"Bytes"`
	Error string   `json:"Error"`
}

// hexBytes decodes space-separated hex, such as "50 44 46".
type hexBytes []byte

func (b *hexBytes) UnmarshalText(text []byte) error {
	decoded, err := hex.DecodeString(strings.ReplaceAll(string(text), " ", ""))
	if err != nil {
		return err
	}
	*b = decoded
	return nil
}

func lookupZXing(t *testing.T) string {
	t.Helper()
	scanner := barcodeScanners(t)["zxing"]
	if !scanner.Enabled {
		t.Skipf("%s: zxing not listed", scannersEnv)
	}
	unavailable := t.Skipf
	if scanner.Required {
		unavailable = t.Fatalf
	}

	path, err := exec.LookPath("ZXingReader")
	if err != nil {
		unavailable("ZXingReader: not found in $PATH")
	}
	out, err := exec.Command(path, "-version").Output()
	if err != nil {
		unavailable("ZXingReader: -version: %v", err)
	}
	version, err := parseZXingVersion(string(out))
	if err != nil {
		unavailable("ZXingReader: -version: %v", err)
	}
	t.Logf("ZXingReader: using %s, version %v", path, version)

	minVersion := scannerVersion{3, 1, 1}
	if compareVersions(version, minVersion) < 0 {
		unavailable("ZXingReader: version %v is older than the minimum %v", version, minVersion)
	}
	return path
}

type scannerOptions struct {
	Enabled  bool
	Required bool
}

type scannerSet map[string]scannerOptions

var scannerUsage = "want auto, all, none or a comma-separated list of: " +
	strings.Join(knownScanners, ", ")

func barcodeScanners(t *testing.T) scannerSet {
	t.Helper()
	v, set := os.LookupEnv(scannersEnv)
	if !set {
		v = "auto"
	}

	scanners := make(scannerSet)
	switch v {
	case "auto", "all":
		for _, name := range knownScanners {
			scanners[name] = scannerOptions{Enabled: true, Required: v == "all"}
		}
		return scanners
	case "none":
		return nil
	}

	for name := range strings.SplitSeq(v, ",") {
		name = strings.TrimSpace(name)
		if len(name) == 0 {
			t.Fatalf("%s: empty scanner name in %q; %s", scannersEnv, v, scannerUsage)
		}
		if !slices.Contains(knownScanners, name) {
			t.Fatalf("%s: unknown scanner %q in %q; %s", scannersEnv, name, v, scannerUsage)
		}
		scanners[name] = scannerOptions{Enabled: true, Required: true}
	}
	return scanners
}

// scannerVersion is a major.minor.patch version.
type scannerVersion [3]int

func (v scannerVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", v[0], v[1], v[2])
}

func parseZXingVersion(s string) (scannerVersion, error) {
	var v scannerVersion
	if _, err := fmt.Sscanf(s, "ZXingReader version %d.%d.%d", &v[0], &v[1], &v[2]); err != nil {
		return v, fmt.Errorf("unexpected output %q: %v", strings.TrimSpace(s), err)
	}
	return v, nil
}

func compareVersions(a, b scannerVersion) int {
	for i := range a {
		if a[i] != b[i] {
			if a[i] < b[i] {
				return -1
			}
			return 1
		}
	}
	return 0
}
