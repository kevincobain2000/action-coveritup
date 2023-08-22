<p align="center">
  <a href="https://coveritup.app">
    <img alt="coveritup app url" src="https://imgur.com/oz9s2zt.png" width="360">
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

# Usages

## Code Coverages

```yaml
    - name: Code Coverage
      run: |
        curl -sLk https://raw.githubusercontent.com/kevincobain2000/cover-totalizer/master/install.sh | sh
        echo SCORE=`./cover-totalizer go_coverage.xml` >> "$GITHUB_ENV"

    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: coverage
    - uses: kevincobain2000/action-coveritup@v1
      with:
        pr_comment: true #if it is a PR, will comment summary
```

## GO Binary Size

```yaml
    - name: Go Binary Size
      run: |
        echo SCORE=`du -sk main | awk '{print $1}'` >> "$GITHUB_ENV"

    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: go-binary-size
```

## GO `go.mod` num of dependencies

```yaml
    - name: Number of dependencies
      run: |
        echo SCORE=`go list -m all|wc -l|awk '{$1=$1};1'` >> "$GITHUB_ENV"

    - uses: kevincobain2000/action-coveritup@v1
      with:
        type: go-mod-dependencies
```

## GO chaining multiple

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

    # Finally comment the summary of 2 reports
    - uses: kevincobain2000/action-coveritup@v1
      with:
        pr_comment: true
```

## Node.js modules size

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

## PHP vendor size

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

# Supported Types

| Type                | Metric |
|---------------------|--------|
| coverage            | %      |
| php-vendor-size     | MB     |
| npm-modules-size    | MB     |
| go-binary-size      | KB     |
| go-mod-dependencies | #      |


---

## Compliance

Destroy just one type

```yaml
      uses: kevincobain2000/action-coveritup@v1
      with:
        destroy: true
        type: npm-modules-size
```

Destroy everything

```yaml
      uses: kevincobain2000/action-coveritup@v1
      with:
        destroy: true
```