import { Router } from 'express';

const router = Router();

/**
 * @openapi
 * /api/votes:
 *   get:
 *     tags:
 *       - Votes
 *     summary: Get list of votes
 *     responses:
 *       200:
 *         description: List of votes
 */
router.get('/', (_req, res) => {
  res.status(200).json({
    success: true,
    data: [],
    message: 'Get votes endpoint - to be implemented',
  });
});

export default router;
