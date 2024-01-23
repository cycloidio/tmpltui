package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	textinput          textinput.Model
	original, rendered textarea.Model
	filepicker         filepicker.Model
	selectedFile       string
	quitting           bool
	err                error

	cwd      string
	tmpldata map[string]any
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m model) Init() tea.Cmd {

	return tea.Batch(tea.EnterAltScreen, m.filepicker.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.original.SetHeight(msg.Height)
		m.original.SetWidth(msg.Width / 2)
		m.rendered.SetHeight(msg.Height)
		m.rendered.SetWidth(msg.Width / 2)
	case clearErrorMsg:
		m.err = nil
	}

	var fcmd, tcmd, icmd tea.Cmd
	m.filepicker, fcmd = m.filepicker.Update(msg)
	m.original, tcmd = m.original.Update(msg)
	m.rendered, tcmd = m.rendered.Update(msg)
	m.textinput, icmd = m.textinput.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		rpath, err := filepath.Rel(m.cwd, path)
		if err != nil {
			// in case of error we just use the abs path
			rpath = path
		}
		m.selectedFile = rpath
		fil, tmpl := loadFile(rpath, m.tmpldata)
		m.original.SetValue(fil)
		m.rendered.SetValue(tmpl)
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		return m, tea.Batch(fcmd, tcmd, icmd, clearErrorAfter(2*time.Second))
	}

	return m, fcmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.selectedFile == "" {
		s.WriteString("Pick a file:")
	} else {
		s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.selectedFile))
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Left, m.textinput.View(), s.String()),
		m.original.View(), m.rendered.View())
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	tea.LogToFile("tmpltui.log", "tmpltui")

	fp := filepicker.New()
	// fp.AllowedTypes = []string{".mod", ".sum", ".go", ".txt", ".md", ".tmpl", ".json", ".yaml"}
	fp.AllowedTypes = []string{}
	fp.CurrentDirectory, _ = os.Getwd()

	m := model{
		filepicker: fp,
		original:   textarea.New(),
		rendered:   textarea.New(),
		cwd:        cwd,
		tmpldata:   loadTmplData(),
	}
	_, err = tea.NewProgram(&m, tea.WithOutput(os.Stderr)).Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("bye")
}

func loadTmplData() map[string]any {
	ret := map[string]any{}
	data, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		log.Fatal(err)
	}

	return ret
}

// loadFile reads the file and attempt to render it as a template, returns
// the original file, the rendered. In case of errors the returned value
// will contain the error text.
func loadFile(fpath string, tdata map[string]any) (string, string) {
	fdata, err := os.ReadFile(fpath)
	if err != nil {
		return fmt.Sprintf("error reading file %s: %+v\n", fpath, err), ""
	}
	t := template.New(fpath).Option("missingkey=zero").Funcs(sprig.FuncMap())

	setDelimiters(t, tdata)

	fstring := string(fdata)
	t, err = t.Parse(fstring)
	if err != nil {
		return fstring, fmt.Sprintf("error parsing template %s: %+v\n", fpath, err)
	}

	var buff bytes.Buffer
	err = t.Execute(&buff, tdata)
	if err != nil {
		return fstring, fmt.Sprintf("error rendering template %s: %+v\n", fpath, err)
	}

	return fstring, buff.String()
}

const (
	leftDelim  = "left_delimiter"
	rightDelim = "right_delimiter"
)

// finds and sets the delimiters
func setDelimiters(t *template.Template, tdata map[string]any) {
	// default delimiters, left and right
	var ld, rd = "{{", "}}"
	if left, ok := tdata[leftDelim]; ok {
		ld = left.(string)
	}
	if right, ok := tdata[rightDelim]; ok {
		rd = right.(string)
	}
	t.Delims(ld, rd)
}
