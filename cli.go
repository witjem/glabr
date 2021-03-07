package main

import (
    "fmt"
    ui "github.com/gizak/termui/v3"
    "github.com/gizak/termui/v3/widgets"
    "github.com/xanzy/go-gitlab"
    "log"
    "os/exec"
    "strings"
    "time"
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
    if err := ui.Init(); err != nil {
       log.Fatalf("failed to initialize termui: %v", err)
    }
    defer ui.Close()

    g := newGui(getAllMrs())
    g.render("")
    // Creating channel using
    // make keyword
    mychan1 := make(chan string, 1)

    // Calling Sleep function of go
    go func() {
       time.Sleep(10 * time.Second)

       // Displayed after sleep overs
       mychan1 <- "output1"
    }()
    events := ui.PollEvents()
    for {
       select {
       case event := <-events:
           if event.ID == "q" || event.ID == "й" {
               return
           }
           if event.ID == "r" || event.ID == "к" {
               g.setMRS(getAllMrs())
           }
           g.render(event.ID)
       case <-mychan1:
           g.setMRS(getAllMrs())
           g.render("")
       }
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
            approved = append(approved, " " + emoji.User.Username)
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

type gui struct {
    mrs         []MR
    owner       *widgets.Paragraph
    list        *widgets.List
    approvals   *widgets.List
    info        *widgets.Paragraph
}

func newGui(mrs []MR) *gui {
    l := widgets.NewList()
    l.Title = "Gitlab MR"
    l.TextStyle = ui.NewStyle(ui.ColorYellow)
    l.WrapText = false
    l.SetRect(1, 1, 100, 30)

    p := widgets.NewParagraph()
    p.Title = "Link"
    p.SetRect(1, 30, 100, 33)

    o := widgets.NewParagraph()
    o.Title = "Owner"
    o.SetRect(100, 1, 120, 4)

    a := widgets.NewList()
    a.Title = "Approvals"
    a.SetRect(100, 3, 120, 33)

    res := &gui{
        mrs:         mrs,
        owner:       o,
        list:        l,
        info:        p,
        approvals:   a,
    }
    res.setMRS(mrs)
    return res
}

func (g *gui) setMRS(mrs []MR) {
    g.list.Rows = []string{}
    for _, mr := range mrs {
        apprv := " "
        if mr.ApprovedByMe {
            apprv = "✓"
        }
        myPr := "  "
        if mr.ApprovedByMe {
            myPr = "my"
        }
        g.list.Rows = append(g.list.Rows, fmt.Sprintf("[%s] [%s] [%s] %s", apprv, myPr, mr.Project, mr.Title))
    }
    g.list.SelectedRow = 0
}

func (g *gui) render(key string) {
    switch key {
    case "q", "<C-c>":
        return
    case "j", "<Down>", "о":
        g.list.ScrollDown()
    case "k", "<Up>", "л":
        g.list.ScrollUp()
    case "o", "щ":
        openURLInBrowser(strings.Trim(g.info.Text, " "))
    }
    g.info.Text = " " + g.mrs[g.list.SelectedRow].Link
    g.owner.Text = " " + g.mrs[g.list.SelectedRow].Owner
    g.approvals.Rows = g.mrs[g.list.SelectedRow].Approved
    ui.Render(g.list, g.info, g.approvals, g.owner)
}

// only works on the macos
func openURLInBrowser(url string) {
    exec.Command("open", url).Start()
}
