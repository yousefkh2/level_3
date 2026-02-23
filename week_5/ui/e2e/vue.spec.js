import { test, expect } from '@playwright/test'

test.describe('PaaS Platform E2E', () => {
  const API_USERNAME = 'admin'
  const API_PASSWORD = 'secret123'

  const uniqueDbName = (prefix = 'test-db') =>
    `${prefix}-${Date.now()}-${Math.floor(Math.random() * 10000)}`

  const dbRow = (page, name) =>
    page
      .locator('li.db-item')
      .filter({ hasText: name })
      .first()

  const loginAsAdmin = async (page) => {
    await page.goto('/')
    await page.fill('input[placeholder="Username"]', API_USERNAME)
    await page.fill('input[placeholder="Password"]', API_PASSWORD)
    await page.click('button[type="submit"]')
    await expect(page.getByRole('heading', { name: 'My Databases' })).toBeVisible()
  }

  const createDatabaseFromUi = async (page, name, instances = '2', storage = '1Gi') => {
    await page.fill('input[placeholder="Database Name"]', name)
    await page.fill('input[placeholder="Instances"]', instances)
    await page.fill('input[placeholder*="Storage"]', storage)
    await page.click('button:has-text("Create")')
    await expect(dbRow(page, name)).toBeVisible({ timeout: 10000 })
  }

  test('complete user flow: login -> create -> view -> delete database', async ({ page }) => {
    const testDbName = uniqueDbName()

    await loginAsAdmin(page)
    await createDatabaseFromUi(page, testDbName)

    const row = dbRow(page, testDbName)
    await row.getByRole('button', { name: 'View Connection Info' }).click()

    const panel = page.locator('.connection-info')
    await expect(panel).toBeVisible()
    await expect(
      panel.getByRole('heading', {
        name: new RegExp(`^Connection Info for ${testDbName}$`),
      })
    ).toBeVisible()
    await expect(panel.getByText(/^Host:/)).toBeVisible()

    await page.click('button:has-text("Close")')

    await row.getByRole('button', { name: 'Delete' }).click()
    await expect(page.locator(`text=${testDbName}`)).not.toBeVisible({ timeout: 10000 })
  })

  test('update database instances and storage via edit flow', async ({ page }) => {
    const testDbName = uniqueDbName('patch-db')

    await loginAsAdmin(page)
    await createDatabaseFromUi(page, testDbName, '2', '1Gi')

    const row = dbRow(page, testDbName)
    await row.getByRole('button', { name: 'Edit' }).click()

    const editForm = row.locator('form.edit-form')
    await expect(editForm).toBeVisible()
    await editForm.locator('input[type="number"]').fill('3')

    const patchResponsePromise = page.waitForResponse((response) => {
      const isPatch = response.request().method() === 'PATCH'
      const isTargetDb = response.url().includes(`/databases/${testDbName}`)
      return isPatch && isTargetDb
    })

    await editForm.getByRole('button', { name: 'Save Changes' }).click()
    const patchResponse = await patchResponsePromise

    expect(patchResponse.status()).toBe(200)

    await expect(row).toContainText('instances: 3', { timeout: 20000 })
    await expect(page.locator('.error')).toHaveCount(0)

    await row.getByRole('button', { name: 'Delete' }).click()
    await expect(dbRow(page, testDbName)).toHaveCount(0, { timeout: 10000 })
  })

  test('login with invalid credentials should fail', async ({ page }) => {
    await page.goto('/')

    await page.fill('input[placeholder="Username"]', 'wrong')
    await page.fill('input[placeholder="Password"]', 'wrong')
    await page.click('button[type="submit"]')

    await expect(page.getByRole('heading', { name: 'Login' })).toBeVisible()
  })
})
