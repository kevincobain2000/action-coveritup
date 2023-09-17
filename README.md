<p align="center">
  <a href="https://coveritup.app">
    <img alt="coveritup app url" src="https://imgur.com/uxGVcRN.png" width="360">
  </a>
</p>
<p align="center">
  Action for uploading Code Metrics to <a href="https://coveritup.app">coveritup.app</a>
</p>

**Quick Setup:** Quickly set up code coverage or other useful metrics on your project.

**Self-Hosted:** Also available. Host your code coverage server.

**Multiple:** Not just code coverage. Track multiple types of reports, such as coverage, lint, bundle size, complexity, etc.

**Pull Request Comments:** Comment on pull requests with the summary report for diff.

**Shield:** Get shields for your `README.md`

**Charts:** Visualize your reports with charts. Report trends over time by branch and user.

# Supported Types

| Type                | Metric | Description                |
|---------------------|--------|----------------------------|
| coverage            | %      |                            |
| php-vendor-size     | MB     |                            |
| npm-modules-size    | MB     |                            |
| go-binary-size      | KB     |                            |
| go-mod-dependencies |        | Num of deps in `go.mod`    |
| go-sec-issues       |        | Go Security issues `gosec` |
| build-time          | s      | Build time in seconds      |


# Embedding shield badges

Shield only

```
![alt title](https://coveritup.app/embed/ORG_NAME/REPO_NAME?branch=BRANCH_NAME&type=TYPE_NAME)
```

Get all your shield badges from https://coveritup.app/ORG_NAME/REPO_NAME

# Using Action

## Code Coverages

```yaml
    - name: Code Coverage
      run: |
        curl -sLk https://raw.githubusercontent.com/kevincobain2000/cover-totalizer/master/install.sh | sh
        echo SCORE=`./cover-totalizer coverage.xml` >> "$GITHUB_ENV"

    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: coverage
    - uses: kevincobain2000/action-coveritup@v1
      with:
        pr_comment: true
```

## Build Time

```yaml
    - name: Build
      run: |
        BUILD_START=$SECONDS
        echo "Your build command. Building..."
        echo SCORE=$(($SECONDS-BUILD_START)) >> "$GITHUB_ENV"

    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: build-time
```

## GO - Binary Size

```yaml
    - name: Go Binary Size
      run: |
        echo SCORE=`du -sk main | awk '{print $1}'` >> "$GITHUB_ENV"

    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: go-binary-size
    - uses: kevincobain2000/action-coveritup@v1
      with:
        pr_comment: true
```

## GO - `go.mod` num of dependencies

```yaml
    - name: Number of dependencies
      run: |
        echo SCORE=`go list -m all|wc -l|awk '{$1=$1};1'` >> "$GITHUB_ENV"

    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: go-mod-dependencies
    - uses: kevincobain2000/action-coveritup@v1
      with:
        pr_comment: true
```

## GO - chaining multiple

```yaml
    # First report
    - name: Go Binary Size
      run: |
        echo SCORE=`du -sk main | awk '{print $1}'` >> "$GITHUB_ENV"
    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: go-binary-size

    # Second report
    - name: go.mod Number of dependencies
      run: |
        echo SCORE=`go list -m all|wc -l|awk '{$1=$1};1'` >> "$GITHUB_ENV"

    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: go-binary-size

    # Finally comment the summary of 2 reports. PR comment should always be in the end.
    - uses: kevincobain2000/action-coveritup@v1
      with:
        pr_comment: true
```

## NPM - modules size

```yaml
    - name: Node Modules Size
      run: |
        echo SCORE=`du -sm node_modules/ | awk '{print $1}'` >> "$GITHUB_ENV"

    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: npm-modules-size
    - uses: kevincobain2000/action-coveritup@v1
      with:
        pr_comment: true
```

## PHP - vendor size

```yaml
    - name: PHP/Composer Vendor Size
      run: |
        echo SCORE=`du -sm vendor/ | awk '{print $1}'` >> "$GITHUB_ENV"

    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: php-vendor-size
    - uses: kevincobain2000/action-coveritup@v1
      with:
        pr_comment: true
```

## PHP - composer number of dependencies

```yaml
    - name: PHP/Composer Vendor Size
      run: |
        echo SCORE=`composer show -i --name-only 2>/dev/null | wc -l | awk '{print $NF}'` >> "$GITHUB_ENV"

    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: composer-dependencies
    - uses: kevincobain2000/action-coveritup@v1
      with:
        pr_comment: true
```

---

## Compliance

**Delete just one type**

```yaml
      uses: kevincobain2000/action-coveritup@v1
      with:
        destroy: true
        type: npm-modules-size
```

**Delete everything**

```yaml
      uses: kevincobain2000/action-coveritup@v1
      with:
        destroy: true
```

**How this action uses `github.token`**

`github.token` from your action is sent to the server as an Authorization header.
The expiration of `github.token` is until the workflow is running.
The token is used to verify if the request has originated from the correct org, repo and commit author.
https://coveritup.app doesn't store the token.
You can see usage in `action.yml` file