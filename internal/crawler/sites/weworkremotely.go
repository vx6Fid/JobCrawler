package sites

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/vx6fid/job-crawler/pkg"
)

type WeWorkRemotelyParser struct{}

func init() {
	Register(&WeWorkRemotelyParser{})
}

func (p *WeWorkRemotelyParser) Matches(url string) bool {
	return strings.Contains(url, "weworkremotely.com")
}

func (p *WeWorkRemotelyParser) Parse(ctx context.Context, e *colly.HTMLElement) ([]pkg.JobPosting, error) {
	log.Println("[STEP] Parse started")

	var jobs []pkg.JobPosting
	e.ForEach("li.new-listing-container.feature > a[href^='/listings/']", func(_ int, el *colly.HTMLElement) {
		select {
		case <-ctx.Done():
			log.Println("--- :| --- Skipping job in Parse() due to timeout.")
			return // this returns from current iteration, not ForEach
		default:
			// continue
		}

		job := pkg.JobPosting{
			Title:    strings.TrimSpace(el.ChildText("h4.new-listing__header__title")),
			Company:  strings.TrimSpace(el.ChildText("p.new-listing__company-name")),
			Location: strings.TrimSpace(el.ChildText("p.new-listing__company-headquarters")),
			ApplyURL: "https://weworkremotely.com" + el.Attr("href"),
		}
		log.Printf("[STEP] -> [weworkremotely] Found job: %s at %s", job.Title, job.Company)
		jobs = append(jobs, job)
	})

	log.Println("[STEP] Parse completed")
	return jobs, nil
}

func (p *WeWorkRemotelyParser) ParseJobDescription(e *colly.HTMLElement) (pkg.JobPosting, error) {
	log.Printf("[STEP] ParseJobDescription started for: %s", e.Request.URL.String())

	job := pkg.JobPosting{
		Title:       strings.ToLower(strings.TrimSpace(e.ChildText("h2.lis-container__header__hero__company-info__title"))),
		Company:     strings.ToLower(strings.TrimSpace(e.ChildText("div.lis-container__header__hero__company-info__description strong"))),
		URL:         e.Request.URL.String(),
		Source:      "weworkremotely.com",
		Description: strings.TrimSpace(e.ChildText("div.lis-container__job__content__description")),
		ApplyURL:    e.ChildAttr("a#job-cta-alt", "href"),
	}

	// Required fields check
	if job.Title == "" || job.Company == "" {
		return pkg.JobPosting{}, fmt.Errorf("failed to parse job description: title or company missing")
	}

	// PostedOn
	postedOnStr := strings.TrimSpace(e.ChildText("li.lis-container__job__sidebar__job-about__list__item:contains('Posted on') span"))
	if postedOnStr != "" {
		parsedTime, err := pkg.ParseRelativeTime(postedOnStr)
		if err == nil {
			job.PostedOn = parsedTime
		} else {
			job.PostedOn = time.Now() // or handle error/log as needed
		}
	}

	// Salary
	job.Salary = strings.TrimSpace(e.ChildText("li.lis-container__job__sidebar__job-about__list__item:contains('Salary') span"))

	// Region/Location (handle absence gracefully)
	location := ""
	e.ForEach("li.lis-container__job__sidebar__job-about__list__item--full:contains('Region') span.box", func(_ int, el *colly.HTMLElement) {
		region := strings.TrimSpace(el.Text)
		if region != "" {
			location += region + ", "
		}
	})
	job.Location = strings.TrimSuffix(location, ", ")

	// Skills from sidebar
	skillsSet := make(map[string]struct{})
	e.ForEach("li.lis-container__job__sidebar__job-about__list__item--full:contains('Skills') span.box", func(_ int, el *colly.HTMLElement) {
		skill := strings.TrimSpace(el.Text)
		if skill != "" {
			skillsSet[skill] = struct{}{}
		}
	})

	// Skills from description (keyword match)
	lowerDesc := strings.ToLower(job.Description)
	for _, keyword := range knownTechnologies {
		if strings.Contains(lowerDesc, strings.ToLower(keyword)) {
			skillsSet[keyword] = struct{}{}
		}
	}

	// Deduplicate and sort skills
	for skill := range skillsSet {
		job.Skills = append(job.Skills, strings.ToLower(skill))
	}
	sort.Strings(job.Skills)

	// Extract experience from description
	job.Experience = extractExperience(job.Description)

	log.Printf("[STEP] -> [weworkremotely] Parsed job: %s at %s", job.Title, job.Company)
	return job, nil
}

func extractExperience(desc string) string {
	re := regexp.MustCompile(`\b(\d{1,2})\+?\s*(years|yrs?)\b`)
	match := re.FindStringSubmatch(strings.ToLower(desc))
	if len(match) > 1 {
		return fmt.Sprintf("%s years", match[1])
	}
	return ""
}
