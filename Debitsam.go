package debitsam

type Debitsam struct {
	eventsam Eventsam
}

func NewDebitsam(eventsam Eventsam) (debitsam *Debitsam, err error) {
	debitsam = &Debitsam{
		eventsam: eventsam,
	}
	return
}
