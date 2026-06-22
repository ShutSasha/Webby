import { test, expect } from '@playwright/test'

test.describe('Admin Control Panel', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')

    await page.getByPlaceholder('Email').fill(process.env.TEST_USER_EMAIL!)
    await page.getByPlaceholder('Password').fill(process.env.TEST_USER_PASSWORD!)

    await page.getByRole('button', { name: 'Log in' }).click()

    await page.waitForURL('/')
  })

  test('should navigate to the admin overview page and display correct title', async ({ page }) => {
    await page.goto('/admin')

    const heading = page.getByRole('heading', { name: 'Admin Control Panel' })
    await expect(heading).toBeVisible()

    const usersStatCard = page.getByText('Total Users')
    await expect(usersStatCard).toBeVisible()
  })

  test('should switch between tabs correctly', async ({ page }) => {
    await page.goto('/admin')

    const usersTab = page.getByRole('link', { name: 'Users Management' })
    await usersTab.click()

    await expect(page).toHaveURL(/.*\/admin\/users/)

    const searchInput = page.getByPlaceholder('Search users by username...')
    await expect(searchInput).toBeVisible()
  })
})
