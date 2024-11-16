# Crafter

> A brief description of your project. Include the purpose, technology stack, and any relevant high-level information.
Crafter is a web-based platform that enables users to easily manage their resumes and make job application and application tracking more simple. It enables users to customize their resumes based on the job descriptions 

## Table of Contents
- [Tech Stack](#tech-stack)
<!-- - [Getting Started](#getting-started) -->
- [Project Structure](#project-structure)
- [Development Guidelines](#development-guidelines)
- [Branching Strategy](#branching-strategy)
- [Pull Request Protocol](#pull-request-protocol)
- [Code Review Process](#code-review-process)

## Tech Stack

- **Frontend:** Next.js (React Framework)
- **Backend:** Go (Golang)
- **Database:** TBD

## Project Structure
/
├── frontend/               # Next.js frontend
├── backend/                # Go backend
└── README.md               # Project documentation

## Developement Guidelines
- Use feature branches for developing new features.
- Commit messages should follow Conventional Commits (e.g., feat: add new authentication flow).

## Branching Strategy
We follow the GitFlow branching strategy:

- main: The stable branch; production-ready code.
- develop: The active development branch.
- Feature branches (feature/your-feature-name): Branches off from develop for adding new features.
- Hotfix branches (hotfix/your-hotfix-name): Directly branched from main to address urgent production issues.

## Pull Request Protocol
- Create a new branch from develop for your feature (feature/your-feature-name).
- Ensure your branch is up to date with the latest changes in develop.
- Before opening a PR, rebase your branch to avoid merge conflicts.
- Add a clear PR title and description summarizing the changes.
- Link any related issue numbers in the PR description.
- Ensure all tests pass before opening a PR.
- PRs should have at least one reviewer (depending on project size, this could be increased).