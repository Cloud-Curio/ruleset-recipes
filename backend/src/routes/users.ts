import { Router } from 'express';

const router = Router();

/**
 * @openapi
 * /api/users/me:
 *   get:
 *     tags:
 *       - Users
 *     summary: Get current user profile
 *     security:
 *       - bearerAuth: []
 *     responses:
 *       200:
 *         description: User profile retrieved successfully
 */
router.get('/me', (_req, res) => {
  res.status(200).json({
    success: true,
    message: 'Get user profile endpoint - to be implemented',
  });
});

export default router;
