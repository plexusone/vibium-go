package vibium_test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/grokify/vibium-go"
)

func Example() {
	ctx := context.Background()

	// Launch browser
	vibe, err := vibium.Launch(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = vibe.Quit(ctx) }()

	// Navigate to a page
	if err := vibe.Go(ctx, "https://example.com"); err != nil {
		log.Fatal(err)
	}

	// Find and click a link
	link, err := vibe.Find(ctx, "a", nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := link.Click(ctx, nil); err != nil {
		log.Fatal(err)
	}

	// Get page title
	title, err := vibe.Title(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Page title:", title)
}

func Example_headless() {
	ctx := context.Background()

	// Launch headless browser
	vibe, err := vibium.LaunchHeadless(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = vibe.Quit(ctx) }()

	// Navigate
	if err := vibe.Go(ctx, "https://example.com"); err != nil {
		log.Fatal(err)
	}

	// Take screenshot
	data, err := vibe.Screenshot(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Save to file
	if err := os.WriteFile("screenshot.png", data, 0644); err != nil {
		log.Fatal(err)
	}
}

func Example_formInteraction() {
	ctx := context.Background()

	vibe, err := vibium.Browser.Launch(ctx, &vibium.LaunchOptions{
		Headless: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = vibe.Quit(ctx) }()

	// Navigate to a form page
	if err := vibe.Go(ctx, "https://example.com/login"); err != nil {
		log.Fatal(err)
	}

	// Fill in username
	username, err := vibe.Find(ctx, "input[name='username']", nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := username.Type(ctx, "myuser", nil); err != nil {
		log.Fatal(err)
	}

	// Fill in password
	password, err := vibe.Find(ctx, "input[name='password']", nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := password.Type(ctx, "mypassword", nil); err != nil {
		log.Fatal(err)
	}

	// Click submit
	submit, err := vibe.Find(ctx, "button[type='submit']", nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := submit.Click(ctx, nil); err != nil {
		log.Fatal(err)
	}
}
