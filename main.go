package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	helpHeight = 5
)

func newTextarea() textarea.Model {
	t := textarea.New()
	t.Prompt = ""
	t.Placeholder = "~these are not the droids you are looking for~"
	t.ShowLineNumbers = true
	t.Cursor.Style = cursorStyle
	t.FocusedStyle.Placeholder = focusedPlaceholderStyle
	t.BlurredStyle.Placeholder = placeholderStyle
	t.FocusedStyle.CursorLine = cursorLineStyle
	t.FocusedStyle.Base = focusedBorderStyle
	t.BlurredStyle.Base = blurredBorderStyle
	t.FocusedStyle.EndOfBuffer = endOfBufferStyle
	t.BlurredStyle.EndOfBuffer = endOfBufferStyle
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.KeyMap.LineNext = key.NewBinding(key.WithKeys("down"))
	t.KeyMap.LinePrevious = key.NewBinding(key.WithKeys("up"))
	t.Blur()
	return t
}

type model struct {
	original, rendered textarea.Model
	filepicker         filepicker.Model
	fpvisible          bool
	selectedFile       string
	quitting           bool
	keymap             keymap
	help               help.Model

	cwd      string
	tmpldata map[string]any
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, m.filepicker.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdList = []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.file):
			m.fpvisible = !m.fpvisible
			cmdList = append(cmdList, func() tea.Msg {
				return focusTexMsg(!m.fpvisible)
			})
		}
	case focusTexMsg:
		if msg {
			m.original.Focus()
			m.rendered.Focus()
		} else {
			m.original.Blur()
			m.rendered.Blur()
		}
	case tea.WindowSizeMsg:
		m.original.SetHeight(msg.Height - helpHeight)
		m.original.SetWidth(msg.Width / 2)
		m.rendered.SetHeight(msg.Height - helpHeight)
		m.rendered.SetWidth(msg.Width / 2)
	}

	var fcmd, tcmd, icmd tea.Cmd
	cmdList = append(cmdList, fcmd, tcmd, icmd)
	if m.fpvisible || m.quitting {
		m.filepicker, fcmd = m.filepicker.Update(msg)
		cmdList = append(cmdList, fcmd)
	}
	m.original, tcmd = m.original.Update(msg)
	m.rendered, tcmd = m.rendered.Update(msg)

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

	return m, tea.Batch(cmdList...)
}

type focusTexMsg bool

func (m model) View() string {
	if m.quitting {
		return ""
	}
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.file,
		m.keymap.quit,
	})
	var views []string
	var s strings.Builder
	if m.fpvisible {
		s.WriteString("\n  ")
		if m.selectedFile == "" {
			s.WriteString("Pick a file:")
		} else {
			s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.selectedFile))
		}
		s.WriteString("\n\n" + m.filepicker.View() + "\n")
	}
	views = append(views, s.String(), m.original.View(), m.rendered.View())
	return lipgloss.JoinHorizontal(lipgloss.Top, views...) + "\n\n" + help
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	tea.LogToFile("tmpltui.log", "tmpltui")

	fp := filepicker.New()
	fp.AllowedTypes = []string{}
	fp.CurrentDirectory, _ = os.Getwd()

	m := model{
		filepicker: fp,
		original:   newTextarea(),
		rendered:   newTextarea(),
		cwd:        cwd,
		tmpldata:   loadTmplData(),
		fpvisible:  true,
		keymap:     newkeymap(),
		help:       help.New(),
	}
	m.original.CharLimit = 8000
	m.rendered.CharLimit = 8000
	_, err = tea.NewProgram(&m, tea.WithOutput(os.Stderr)).Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("bye")
}
