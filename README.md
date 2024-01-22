<p align="center">
  <a href="https://coveritup.app">
    <img alt="coveritup app url" src="https://imgur.com/fHfULta.png" width="160">
  </a>
</p>

<p align="center">
  <a href="https://coveritup.app">
    <img alt="coveritup app url" src="https://imgur.com/7pQEwvT.png" width="460">
  </a>
</p>

<p align="center">
  Action for uploading Code Metrics to <a href="https://coveritup.app">coveritup.app</a>
</p>

![npm-install-time](https://coveritup.app/badge?org=kevincobain2000&repo=action-coveritup&type=npm-install-time&branch=master)
![npm-build-time](https://coveritup.app/badge?org=kevincobain2000&repo=action-coveritup&type=npm-build-time&branch=master)
![go-build-time](https://coveritup.app/badge?org=kevincobain2000&repo=action-coveritup&type=go-build-time&branch=master)

![coverage](https://coveritup.app/badge?org=kevincobain2000&repo=action-coveritup&type=coverage&branch=master)
![go-binary-size](https://coveritup.app/badge?org=kevincobain2000&repo=action-coveritup&type=go-binary-size&branch=master)

![go-test-run-time](https://coveritup.app/badge?org=kevincobain2000&repo=action-coveritup&type=go-test-run-time&branch=master)
![go-mod-dependencies](https://coveritup.app/badge?org=kevincobain2000&repo=action-coveritup&type=go-mod-dependencies&branch=master)
![go-sec-issues](https://coveritup.app/badge?org=kevincobain2000&repo=action-coveritup&type=go-sec-issues&branch=master)

**Quick Setup:** Quickly set up code coverage or other useful metrics on your project.

**Self-Hosted:** Also available. Host your code coverage server.

**Multiple:** Not just code coverage. Track multiple types of reports, such as coverage, lint, bundle size, complexity, etc.

**Pull Request Comments:** Comment on pull requests with the summary report for diff.

**Shield:** Get shields for your `README.md`

**Charts:** Visualize your reports with charts. Report trends over time by branch and user.

- [Step 1) Using Action](#step-1-using-action)
- [Step 2) Add to your workflow](#step-2-add-to-your-workflow)
  - [Scores](#scores)
  - [Times](#times)
  - [Sizes](#sizes)
  - [Counts](#counts)
- [Step 3) Embedding shield badges in README](#step-3-embedding-shield-badges-in-readme)
- [Compliance](#compliance)


# Step 1) Using Action

Before using this action, enable Github Actions

![Github Settings](https://imgur.com/psKpD15.png)

- [x] Read and write permission
- [x] Allow Github Actions to create and approve pull requests


# Step 2) Add to your workflow

## Scores

```yaml
    # Example: Clover
    - name: Code Coverage
      run: |
        curl -sLk https://raw.githubusercontent.com/kevincobain2000/cover-totalizer/master/install.sh | sh
        echo SCORE=`./cover-totalizer coverage.xml` >> "$GITHUB_ENV"
    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: coverage

    # Finally comment on PR
    - uses: kevincobain2000/action-coveritup@v1
      with:
        pr_comment: true
```

## Times

```yaml
    # Example: Go
    - name: Build
      run: |
        BUILD_START=$SECONDS
        go build main.go
        echo SCORE=$(($SECONDS-BUILD_START)) >> "$GITHUB_ENV"
    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: go-build-time

    # Example: NPM
    - name: Build
      run: |
        BUILD_START=$SECONDS
        npm install
        echo SCORE=$(($SECONDS-BUILD_START)) >> "$GITHUB_ENV"
    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: npm-build-time

    # Finally comment on PR
    - uses: kevincobain2000/action-coveritup@v1
      with:
        pr_comment: true
```

## Sizes

```yaml
    # Example: Go
    - name: Go Binary Size
      run: |
        echo SCORE=`du -sk main | awk '{print $1}'` >> "$GITHUB_ENV"
    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: go-binary-size

    # Example: NPM
    - name: Node Modules Size
      run: |
        echo SCORE=`du -sm node_modules/ | awk '{print $1}'` >> "$GITHUB_ENV"
    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: npm-modules-size

    # Example: PHP
    - name: PHP/Composer Vendor Size
      run: |
        echo SCORE=`du -sm vendor/ | awk '{print $1}'` >> "$GITHUB_ENV"
    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: php-vendor-size

    # Finally comment on PR
    - uses: kevincobain2000/action-coveritup@v1
      with:
        pr_comment: true
```

## Counts

```yaml
    # Example: Go
    - name: Number of dependencies
      run: |
        echo SCORE=`go list -m all|wc -l|awk '{$1=$1};1'` >> "$GITHUB_ENV"
    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: go-mod-dependencies

    # Example: PHP
    - name: PHP/Composer Vendor Size
      run: |
        echo SCORE=`composer show -i --name-only 2>/dev/null | wc -l | awk '{print $NF}'` >> "$GITHUB_ENV"
    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: composer-dependencies

    # Finally comment on PR
    - uses: kevincobain2000/action-coveritup@v1
      with:
        pr_comment: true
```

# Step 3) Embedding shield badges in README

Navigate to your repo and obtain shield badges https://coveritup.app/explore

---

# Compliance

**Delete just one type**

```yaml
    - uses: kevincobain2000/action-coveritup@v1
      with:
        destroy: true
        type: npm-modules-size
```

**Delete everything**

```yaml
    - uses: kevincobain2000/action-coveritup@v1
      with:
        destroy: true
```

**How this action uses `github.token`**

`github.token` from your action is sent to the server as an Authorization header.
The expiration of `github.token` is until the workflow is running.
The token is used to verify if the request has originated from the correct org, repo and commit author.
https://coveritup.app doesn't store the token.
You can see usage in `action.yml` file