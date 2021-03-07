package main

import (
    "fmt"
    "github.com/gizak/termui/v3"
    "github.com/gizak/termui/v3/widgets"
    "github.com/xanzy/go-gitlab"
    "log"
    "os/exec"
    "syscall"
    "time"
    "unsafe"
)

type MR struct {
    Link         string
    Project      string
    Title        string
    ApprovedByMe bool
    Owner        string
    Approved     []string
}

func main() {
    if err := termui.Init(); err != nil {
        log.Fatalf("failed to initialize termui: %v", err)
    }
    defer termui.Close()

    ui := newUI()
    ui.render()

    gitlabUpdates := make(chan []MR, 1)
    go func() {
        gitlabUpdates <- getAllMrs()
        time.Sleep(10 * time.Second)
    }()
    events := termui.PollEvents()
    for {
        select {
        case event := <-events:
            switch event.ID {
            case "q", "й":
                return
            case "j", "<Down>", "о":
                ui.ScrollDown()
            case "k", "<Up>", "л":
                ui.ScrollUp()
            case "o", "щ":
                openURLInBrowser(ui.CurrentMR().Link)
            }
        case mrs := <-gitlabUpdates:
            ui.updateMrs(mrs)
        }
        ui.render()
    }
}

func getAllMrs() []MR {
    mrsM := [][]MR{
        getMrs(243, "ims"),
        getMrs(245, "fle"),
        getMrs(277, "dss"),
        getMrs(297, "ats"),
        getMrs(246, "gtw"),
    }
    mrs := []MR{}
    for _, mr := range mrsM {
        mrs = append(mrs, mr...)
    }
    return mrs
}

func getMrs(pid int, pName string) []MR {
    git, err := gitlab.NewClient("k9zB-yyfBSS9tYTyYmxW", gitlab.WithBaseURL("https://gitlab.railsreactor.com/"))
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    mrs, _, _ := git.MergeRequests.ListProjectMergeRequests(pid, &gitlab.ListProjectMergeRequestsOptions{
        ListOptions: gitlab.ListOptions{
            Page:    0,
            PerPage: 10,
        },
        State: gitlab.String("opened"),
    })
    res := []MR{}
    for _, mr := range mrs {
        emojis, _, _ := git.AwardEmoji.ListMergeRequestAwardEmoji(pid, mr.IID, &gitlab.ListAwardEmojiOptions{PerPage: 10})
        isApproved := false
        approved := []string{}
        for _, emoji := range emojis {
            if emoji.User.Username == "roman.ilnytskyi" {
                isApproved = true
            }
            approved = append(approved, " "+emoji.User.Username)
        }
        res = append(res, MR{
            Link:         mr.WebURL,
            Project:      pName,
            Title:        mr.Title,
            ApprovedByMe: isApproved,
            Owner:        mr.Author.Username,
            Approved:     approved,
        })
    }
    return res
}

type ui struct {
    mrs           []MR
    ownerView     *widgets.Paragraph
    listView      *widgets.List
    approvalsView *widgets.List
    linkView      *widgets.Paragraph
}

func newUI() *ui {
    maxX, maxY := termSize()

    listView := widgets.NewList()
    listView.Title = "Gitlab MR"
    listView.TextStyle = termui.NewStyle(termui.ColorYellow)
    innerX := int(float32(maxX) * 0.7)
    listView.SetRect(1, 1, innerX, maxY-3)

    linkView := widgets.NewParagraph()
    linkView.Title = "Link"
    linkView.SetRect(1, maxY-3, innerX, maxY)

    ownerView := widgets.NewParagraph()
    ownerView.Title = "Owner"
    ownerView.SetRect(innerX, 1, maxX, 4)

    approvalsView := widgets.NewList()
    approvalsView.Title = "Approvals"
    approvalsView.SetRect(innerX, 4, maxX, maxY)

    res := &ui{
        mrs: []MR{
            {
                Link:         "",
                Project:      "",
                Title:        "",
                ApprovedByMe: false,
                Owner:        "",
                Approved:     []string{},
            },
        },
        ownerView:     ownerView,
        listView:      listView,
        linkView:      linkView,
        approvalsView: approvalsView,
    }
    return res
}

func (u *ui) updateMrs(mrs []MR) {
    u.listView.Rows = []string{}
    for _, mr := range mrs {
        approveSymbol := " "
        if mr.ApprovedByMe {
            approveSymbol = "✓"
        }
        myMrSymbol := "  "
        if mr.ApprovedByMe {
            myMrSymbol = "my"
        }
        u.listView.Rows = append(u.listView.Rows, fmt.Sprintf("[%s] [%s] [%s] %s", approveSymbol, myMrSymbol, mr.Project, mr.Title))
    }
    u.mrs = mrs
    u.listView.SelectedRow = 0
}

func (u *ui) render() {
    u.linkView.Text = " " + u.mrs[u.listView.SelectedRow].Link
    u.ownerView.Text = " " + u.mrs[u.listView.SelectedRow].Owner
    u.approvalsView.Rows = u.mrs[u.listView.SelectedRow].Approved
    termui.Render(u.listView, u.linkView, u.approvalsView, u.ownerView)
}

func (u *ui) ScrollDown() {
    u.listView.ScrollDown()
}

func (u *ui) ScrollUp() {
    u.listView.ScrollUp()
}

func (u *ui) CurrentMR() MR {
    return u.mrs[u.listView.SelectedRow]
}

// only works on the macos
func openURLInBrowser(url string) {
    exec.Command("open", url).Start()
}

// only works on the macos, and linux
func termSize() (int, int) {
    var sz struct {
        rows    uint16
        cols    uint16
        xpixels uint16
        ypixels uint16
    }
    _, _, _ = syscall.Syscall(syscall.SYS_IOCTL,
        uintptr(syscall.Stdout), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&sz)))
    return int(sz.cols), int(sz.rows)
}
