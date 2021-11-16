import { markdown, message, warn, danger } from 'danger'
import { getTestCoverageEntries } from './utils'

const config = {
  tests: {
    artifactPath: process.env.TEST_COVERAGE_ARTIFACT_PATH ?? '../artifacts/coverage.out',
    threshold: +(process.env.TEST_COVERAGE_THRESHOLD ?? 80)
  },
  linesOfCodeThreshold: +(process.env.LINES_OF_CODE_THRESHOLD ?? 500),
  minimumPrDescriptionLength: +(process.env.MINIMUM_PR_DESCRIPTION_LENGTH ?? 50)
}

const testEntries = getTestCoverageEntries(config.tests.threshold, config.tests.artifactPath)

// Markdown overall level of testing coverge
markdown(`Overall test coverage level is **${testEntries.pop().coverage}%** and the threshold is ${config.tests.threshold}%.`)

// Warn if there is at least one test coverage below the threshold.
const warnings = testEntries.filter((entry) => entry.isBelowThreshold)
if (warnings.length) {
  warn(`There are some functions(s) below ${config.tests.threshold}% threshold.`)
  const header = '|Coverage|Line Number|Function|\n|:-:|:--|:--|\n'
  const rows = warnings
    .map(({ fileName, method, coverage }) => `|${coverage}%|\`${fileName}\`|\`${method}\`|`)
    .join('\n')
  markdown(`${header}${rows}`)
}

// Warn if the pull request is big.
(async () => {
  const lines = await danger.git.linesOfCode()
  if (lines > config.linesOfCodeThreshold) {
    warn(`Big PR! ${lines} lines of code...`)
  }
})()

// Warn if `go.mod` is being updated without a change in `go.sum`.
const allChangedFiles =
  danger.git.created_files
    .concat(danger.git.modified_files)
    .concat(danger.git.deleted_files)

if (allChangedFiles.includes('go.mod') && !allChangedFiles.includes('go.sum')) {
  warn('`go.mod` has been updated without a change in `go.sum`.')
}

// Warn for a short or empty PR body.
if (danger.github.pr.body == null || danger.github.pr.body.trim().length === 0) {
  warn('The PR description is missing!')
} else if (danger.github.pr.body.trim().length < config.minimumPrDescriptionLength) {
  warn(`The PR description is too short! Please describe a little bit more than ${config.minimumPrDescriptionLength} chars.`)
}
