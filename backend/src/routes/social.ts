import { Router } from 'express';

const router = Router();

/**
 * @openapi
 * /api/social/posts:
 *   get:
 *     tags:
 *       - Social
 *     summary: Get social feed posts
 *     security:
 *       - bearerAuth: []
 *     responses:
 *       200:
 *         description: List of posts
 */
router.get('/posts', (_req, res) => {
  res.status(200).json({
    success: true,
    data: [],
    message: 'Get social posts endpoint - to be implemented',
  });
});

export default router;
