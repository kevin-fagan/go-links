<p align="center">
  <img src="https://github.com/user-attachments/assets/beffcad9-83c8-43a0-859c-af0eadb22150" alt="golinks-icon" width="150" />
</p>

Go Links is a lightweight, self-hosted application for creating easy-to-remember vanity URLs. Great for internal tools, sites or personal projects.
Instead of sharing or having to rember long, complex URLs, you can create short aliases such as `go/standup` or `go/jira`.

- ⚡ Fast and light weight 
- 📝 Built in audit trail. Every change is recorded for traceability
- 🔒 Fully self hosted. You maintain complete control over your data and links
- 🔋 Batteries included. No 3rd party network calls for HTML, JS, Fonts or CSS


| Links                                  | Logs                                   |
| -------------------------------------- | -------------------------------------- |
| ![link-page](./assets/links-page.png)  | ![logs-page](./assets/logs-page.png)   |

## Motivation

Go Links is not a new concept. It's been used within tech companies for years to simplify internal URL sharing. But, there doesn’t seem to be a great open-source version available. Most options are either poorly maintained or are locked behind a paywall. And why should someone have to pay for something as simple as a URL redirect? So I decided to make my own "Go Links" with hopes that others may find it useful.  

## Setup

Docker and configuration files are coming soon. Still, getting the project to run locally is easy!
1. `make setup`
2. `make run`

## Roadmap

Additional planned features include:

- 🏷️ Tagging system
- 📤 Ability to export your data 
- 🔐 Authentication / Permissions
