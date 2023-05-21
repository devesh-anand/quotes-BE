package cron

import (
	"log"
	"time"

	"api.deveshanand.com/quotes"
	"api.deveshanand.com/quotes/types"

	"github.com/robfig/cron"
)

var cronQuotes []types.PostData = []types.PostData{
	{Quote: "All men make mistakes, but a good man yields when he knows his course is wrong, and repairs the evil. The only crime is pride.", Author: "Sophocles", Sub_by: "Devesh Anand"},
	{Quote: "He who has a why to live can bear almost any how.", Author: "Friedrich Nietzsche", Sub_by: "Devesh Anand"},
	{Quote: "Unexpressed emotions will never die. They are buried alive and will come forth later in uglier ways.", Author: "Sigmund Freud", Sub_by: "Devesh Anand"},
	{Quote: "No tree can grow to heaven unless its roots reach down to hell.", Author: "Carl Jung", Sub_by: "Devesh Anand"},
	{Quote: "If you wish to improve, be content to appear clueless or stupid.", Author: "Epictetus", Sub_by: "Devesh Anand"},
	{Quote: "You have no enemies. No one has any enemies. There is no one that you should hurt.", Author: "Thors Snorresson (Vinland Saga)", Sub_by: "Devesh Anand"},
	{Quote: "If a man know not to which port he sails, no wind is favourable.", Author: "Seneca", Sub_by: "Devesh Anand"},
	{Quote: `And God said, "Love your enemy", and I obeyed him and loved myself.`, Author: "Kahlil Gibran", Sub_by: "Devesh Anand"},
	{Quote: "You are in the danger of living a life so comfortable and soft, that you will die without ever realising your true potential.", Author: "David Goggins", Sub_by: "Devesh Anand"},
	{Quote: "Morallty Is Just a fiction used by the herd of inferior human beings to hold back the few superior men.", Author: "Fredrich Nietzsche", Sub_by: "Devesh Anand"},
	{Quote: "Misfortune tests great men.", Author: "Seneca", Sub_by: "Devesh Anand"},
	{Quote: "It never ceases to amaze me: we love ourselves more than other people, but care more about their opinions than our own.", Author: "Marcus Aurelius", Sub_by: "Devesh Anand"},
	{Quote: "Leisure without study is death - a tomb for the living.", Author: "Seneca", Sub_by: "Devesh Anand"},
}
var cronNum int = 0
func SubmitCron() {
	c := cron.New();
	c.AddFunc("* * * */7 * *", func() {
		if(cronNum < len(cronQuotes)) {
			log.Println(cronQuotes[cronNum], time.Now())

			quotes.AddQuote(cronQuotes[cronNum], 1)
			cronNum++
		}
	})

	c.Start()
}