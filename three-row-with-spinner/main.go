package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	windowHeight = 0
	windowWidth  = 0
)

type model struct {
	spinner spinner.Model
	busy    bool
}

func main() {
	var m model

	s := initialSpinnerModel()
	m.busy = true
	m.spinner = s
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// if your program needs to do something as soon as it's loaded
// do it here
func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, mockSpendTime)
}

// classic Update() function that all bubbletea apps need
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	// This means the cmd goroutine has returned
	// so we will never show the spinner again
	case mockSpendTimeMsg:
		m.busy = false
		return m, cmd

	case tea.WindowSizeMsg:
		windowHeight = msg.Height
		windowWidth = msg.Width
		return m, cmd

	case tea.KeyMsg:

		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	default:
		// we have to update the spinner for it to actually animate.
		// This is a convenient place to do it.
		if m.busy {
			m.spinner, cmd = m.spinner.Update(msg)
		}
		return m, cmd
	}
	return m, cmd
}

// classic View() function that all bubbletea apps need
func (m model) View() string {

	/* usually people put all this stuff in functions. Left here for simplicity sake */
	
	// note that the horizontal lines are nothing but cleverly placed borders on just top or bottom of all elements
	// I don't know if there's a way to use a container element and not have to do this for each section
	topBoxLeft := lipgloss.NewStyle().Width(fixSize(windowWidth, 2)).AlignHorizontal(lipgloss.Left).BorderStyle(lipgloss.NormalBorder()).BorderBottom(true)
	topBoxRight := lipgloss.NewStyle().Width(windowWidth / 2).AlignHorizontal(lipgloss.Right).BorderStyle(lipgloss.NormalBorder()).BorderBottom(true)

	topBoxContent := lipgloss.JoinHorizontal(lipgloss.Left,
		topBoxLeft.Render("Top Left"),
		topBoxRight.Render("Top Right"))

	bottomBoxLeft := lipgloss.NewStyle().Width(windowWidth / 3).AlignHorizontal(lipgloss.Left).BorderStyle(lipgloss.NormalBorder()).BorderTop(true)
	bottomBoxMiddle := lipgloss.NewStyle().Width(windowWidth / 3).AlignHorizontal(lipgloss.Center).BorderStyle(lipgloss.NormalBorder()).BorderTop(true)
	bottomBoxRight := lipgloss.NewStyle().Width(fixSize(windowWidth, 3)).AlignHorizontal(lipgloss.Right).BorderStyle(lipgloss.NormalBorder()).BorderTop(true)

	bottomBoxContent := lipgloss.JoinHorizontal(lipgloss.Bottom,
		bottomBoxLeft.Render("Bottom Left "),
		bottomBoxMiddle.Render("Bottom Center"),
		bottomBoxRight.Render("Bottom Right"))

	// this is here because we need the final sizes of the top and bottom bars before we can assign a dynamic size to the middle
	middleBox := lipgloss.NewStyle().Width(windowWidth).Height(windowHeight - 4)
	// this 4 above is my top bar height of 2 + bottom bar height of 2 . could be set and queried instead of hardcoded. (I can't because I didn't set them)

	middleBoxContent := ""

	// this is just some size info to display in the middle box. not important how it's put together.
	content := fmt.Sprintf("\n\nWindow Height: %d Width: %d", windowHeight, windowWidth)
	content += fmt.Sprintf("\n\nHeight (middleBox: %d) ", middleBox.GetHeight())
	content += fmt.Sprintf("\n\nWidth (topBox: %d , middleBox:, %d bottomBox: %d ) ", topBoxLeft.GetWidth(), middleBox.GetWidth(), bottomBoxLeft.GetWidth())

	// now that top/bottoms sections are ready we can display something
	if m.busy {
		middleBoxContent = middleBox.AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render(m.spinner.View() + " A lot of important work is being done! " + m.spinner.View())
	} else {
		middleBoxContent = middleBox.Width(windowWidth).AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render("Resize me to see window/box sizes change while the top and bottom boxes remain static!" + content)
	}
	return lipgloss.JoinVertical(lipgloss.Left, topBoxContent, middleBoxContent, bottomBoxContent)
}

// we need this just  so we can return a type tea.Msg from mockSpendTimeMsg()
// just using bool since tea.Msg can be any type
type mockSpendTimeMsg bool

// normally this function would do some actual work
// we sleep here because otherwise the spinner will
// disappear before we can see it. In real life if work
// completes that quickly, we'd like it even more.
func mockSpendTime() tea.Msg {
	time.Sleep(2 * time.Second)

	// note that you will catch this inside the Update() loop
	return mockSpendTimeMsg(true)
}

// create a spinner component with a nice color
func initialSpinnerModel() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return s
}

// because window width/height are integers
// it is not always possible for total size of sections
// to exactly match the total window sizes (i.e 100/3 = 33 while 3*33=99)
// This function when called in the right section will add the missing little
// bit to one of the sections so things appear whole
func fixSize(total int, parts int) int {
	remainder := total % parts
	if remainder > 0 {
		return (total / parts) + remainder
	} else {
		return total / parts
	}
}
