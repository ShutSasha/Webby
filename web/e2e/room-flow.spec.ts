import { test, expect } from '@playwright/test'

test.describe('Room End-to-End Flow', () => {
  const uniqueRoomName = `E2E Test Room ${Date.now()}`
  const testMessage = 'Hello from Playwright E2E test!'

  test('User can create room, add video, message to chat, and delete room', async ({ page }) => {
    await test.step('Login to the platform', async () => {
      await page.goto('/login')
      await page.getByPlaceholder('Email').fill(process.env.TEST_USER_EMAIL!)
      await page.getByPlaceholder('Password').fill(process.env.TEST_USER_PASSWORD!)
      await page.getByRole('button', { name: 'Log in' }).click()
      await page.waitForURL('**/')
    })

    await test.step('Create a new room and enter it', async () => {
      await page.goto('/rooms')

      await page.getByRole('button', { name: /create room/i }).click()
      await page.getByPlaceholder(/Enter room name.../i).fill(uniqueRoomName)
      await page.locator('select').selectOption({ index: 1 })
      await page.locator('form').getByRole('button', { name: 'Create Room' }).click()
      await expect(page.getByRole('heading', { name: 'Create a new room' })).not.toBeVisible()

      await page.getByRole('link', { name: /my rooms/i }).click()
      const newRoomCard = page.locator('div').filter({ hasText: uniqueRoomName }).first()
      await expect(newRoomCard).toBeVisible()
      await newRoomCard.getByRole('link', { name: `Enter ${uniqueRoomName}` }).click()

      await expect(page.getByPlaceholder('Search video in queue')).toBeVisible()
      await expect(page.getByText(uniqueRoomName)).toBeVisible()
    })

    await test.step('Add a video from Webby platform', async () => {
      await page.getByRole('button', { name: 'Search' }).click()

      const searchInput = page.getByPlaceholder(/Search users, rooms/i)
      await expect(searchInput).toBeVisible()

      await searchInput.fill('a')

      const firstVideoCard = page.getByRole('link').filter({ hasText: /views/i }).first()
      await expect(firstVideoCard).toBeVisible({ timeout: 10000 })

      await firstVideoCard.hover()

      const plusButton = firstVideoCard.locator('button').last()
      await plusButton.click()

      await page.getByRole('button', { name: /add to room/i }).click()
    })

    await test.step('Switch to chat and send a message', async () => {
      await page.keyboard.press('Escape')
      await expect(page.getByPlaceholder(/Search users, rooms/i)).not.toBeVisible()

      await page.getByTestId('tab-chat').click()

      await expect(page.getByText('No messages yet. Be the first to say hello!')).toBeVisible({ timeout: 10000 })

      const chatInput = page.getByPlaceholder('Send a message')
      await expect(chatInput).toBeVisible()

      await chatInput.fill(testMessage)
      await chatInput.press('Enter')

      await expect(page.getByText(testMessage)).toBeVisible()
    })

    await test.step('Navigate to own rooms and delete the room', async () => {
      await page.goto('/rooms')

      await page.getByRole('link', { name: /my rooms/i }).click()

      const roomCard = page.locator('.group').filter({
        has: page.getByRole('heading', { name: uniqueRoomName, exact: true }),
      })
      await expect(roomCard).toBeVisible()

      await roomCard.hover()
      await roomCard.locator('button').click()

      await page.getByRole('button', { name: 'Delete room' }).click()

      await expect(page.getByText(uniqueRoomName)).not.toBeVisible()
    })
  })
})
