<p align="center">
  <img loading='eager' alt="coveritup app url" src="https://imgur.com/fHfULta.png" width="160">
</p>
<p align="center">
  The Hassle-Free Code Coverage and Metrics Tracking Tool
  <br>
  A CodeCov and SonarQube alternative.
  <br>
</p>

<p align="center">
  <img loading='eager' alt="coveritup app url" src="https://imgur.com/7pQEwvT.png" width="460">
</p>

<p align="center">
  Action and Self Hosted app for uploading Code Metrics.
</p>

<p align="center">
    <a
        href="https://www.producthunt.com/posts/coveritup?utm_source=badge-featured&utm_medium=badge&utm_souce=badge-coveritup"
        target="_blank"
        ><img loading='eager'
        src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=433114&theme=light"
        alt="CoverItUp - All&#0032;in&#0032;one&#0032;code&#0032;coverage&#0032;and&#0032;badges&#0032;tool&#0046; | Product Hunt"
        style="width: 250px; height: 54px;"
        width="250"
        height="54"
        /></a
    >
</p>

---

<br>
<p align="center">
    <b>Add Progress To your README</b>
    <br>
    <br>
    <img loading='eager' src="https://raw.githubusercontent.com/kevincobain2000/action-coveritup/refs/heads/master/samples/1.svg" /> &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
    <img loading='eager' src="https://raw.githubusercontent.com/kevincobain2000/action-coveritup/refs/heads/master/samples/2.svg" />
</p>

---

<p align="center">
    <b>Add Badges To your README</b>
</p>

![npm-install-time](https://raw.githubusercontent.com/kevincobain2000/action-coveritup/refs/heads/master/samples/3.svg)

![coverage](https://raw.githubusercontent.com/kevincobain2000/action-coveritup/refs/heads/master/samples/4.svg)

---

<p align="center">
    <b>Add Widgets To your README</b>
</p>

![npm-install-time](https://raw.githubusercontent.com/kevincobain2000/action-coveritup/refs/heads/master/samples/5.svg)
![go-test-run-time](https://raw.githubusercontent.com/kevincobain2000/action-coveritup/refs/heads/master/samples/6.svg)

---

**Quick Setup:** Quickly set up code coverage or other useful metrics on your project.

**Self-Hosted:** Also available. Host your code coverage server.

**Multiple:** Not just code coverage. Track multiple types of reports, such as coverage, lint, bundle size, complexity, etc.

**Pull Request Comments:** Comment on pull requests with the summary report for diff.

**Shield:** Get shields for your `README.md`

**Charts:** Visualize your reports with charts. Report trends over time by branch and user.

**Read on Medium** - [Revolutionizing Code Tracking for Developers](https://kevincobain2000-x.medium.com/revolutionizing-code-tracking-for-developers-e8b7b42a5204), [Use Linters the Right way](https://kevincobain2000-x.medium.com/revolutionizing-code-tracking-for-developers-e8b7b42a5204)


# Examples - Pull Request Comments

Do a CF (Continuous Feedback) on your pull requests. Comment on PR with the summary report for diff.

See this pull request for example: https://github.com/kevincobain2000/action-coveritup/pull/15


# Step 1) Using Action

Before using this action, enable Github Actions

![Github Settings](https://imgur.com/psKpD15.png)

- [x] Read and write permission
- [x] Allow Github Actions to create and approve pull requests


# Step 2) Add to your workflow

## Scores `example code coverage`

```yaml
    # Example: Clover
    - run: curl -sLk https://raw.githubusercontent.com/kevincobain2000/cover-totalizer/master/install.sh | sh
    - uses: kevincobain2000/action-coveritup@v2
      with:
        type: coverage
        command: ./cover-totalizer coverage.xml

    # Finally comment on PR
    - uses: kevincobain2000/action-coveritup@v2
      with:
        pr_comment: true
        # optional
        ## report only these types on PR comment, empty means all
        types: coverage,go-sec-issues,go-lint-errors
        # optional
        ## report only these types after 1st comment
        # 1st comment will have all types or types specified in `types`
        # 2nd comment onwards will have only these types
        diff_types: coverage
```

## Time taken

```yaml
    # Example: Go
    - uses: kevincobain2000/action-coveritup@v2
      with:
        type: go-build-time
        command: go build main.go
        record: runtime

    # Example: NPM
    - uses: kevincobain2000/action-coveritup@v2
      with:
        type: npm-build-time
        command: npm run build
        record: runtime

    # Finally comment on PR
    - uses: kevincobain2000/action-coveritup@v2
      with:
        pr_comment: true
```

## Bundle sizes

```yaml
    # Example: Go
    - uses: kevincobain2000/action-coveritup@v2
      with:
        type: go-binary-size
        command: du -sk main | awk '{print $1}'

    # Example: NPM
    - uses: kevincobain2000/action-coveritup@v2
      with:
        type: npm-modules-size
        command: du -sm node_modules/ | awk '{print $1}'

    # Example: PHP
    - uses: kevincobain2000/action-coveritup@v2
      with:
        type: php-vendor-size
        command: du -sm vendor/ | awk '{print $1}'

    # Finally comment on PR
    - uses: kevincobain2000/action-coveritup@v2
      with:
        pr_comment: true
```

## Counts

```yaml
    # Example: Go
    - uses: kevincobain2000/action-coveritup@v2
      with:
        type: go-mod-dependencies
        command: go list -m all|wc -l|awk '{$1=$1};1'

    # Example: PHP
    - uses: kevincobain2000/action-coveritup@v2
      with:
        type: composer-dependencies
        command: composer show -i --name-only 2>/dev/null | wc -l | awk '{print $NF}'

    # Finally comment on PR
    - uses: kevincobain2000/action-coveritup@v2
      with:
        pr_comment: true
```

# Step 3) Embedding badges and charts

Navigate to your repo and obtain embeding code for badges and charts.
`/readme?org=kevincobain2000&repo=action-coveritup&branch=master`

---

# Compliance

**Delete just one type**

```yaml
    - uses: kevincobain2000/action-coveritup@v2
      with:
        destroy: true
        type: npm-modules-size
```

**Delete everything**

```yaml
    - uses: kevincobain2000/action-coveritup@v2
      with:
        destroy: true
```

**How this action uses `github.token`**

`github.token` from your action is sent to the server as an Authorization header.
The expiration of `github.token` is until the workflow is running.
The token is used to verify if the request has originated from the correct org, repo and commit author.
It doesn't store the token.
You can see usage in `action.yml` file


# Development notes

```sh
# for backend
cd app/
air # or go run main.go

# for frontend
cd app/fronend
npm install
npm run dev
```


# Self Hosting Options

## Just the Api

Download the binary from [releases](https://github.com/kevincobain2000/action-coveritup/releases)

## Build from source (with UI)

```sh
git clone https://github.com/kevincobain2000/action-coveritup
cd app
cp frontend/.env.example frontend/.env
./build.sh
./main
```

# CHANGE LOG

- **v1.0** - Initial release with `self hosted` and `action`.
- **v2.0** - Better action that wraps the command.
- **v2.4** - Smoothed bar charts.
