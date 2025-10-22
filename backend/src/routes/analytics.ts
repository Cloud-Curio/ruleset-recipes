import { Router } from 'express';

const router = Router();

/**
 * @openapi
 * /api/analytics/politicians/{id}/similarity:
 *   get:
 *     tags:
 *       - Analytics
 *     summary: Get politician similarity scores
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Similarity scores
 */
router.get('/politicians/:id/similarity', (_req, res) => {
  res.status(200).json({
    success: true,
    data: [],
    message: 'Get politician similarity endpoint - to be implemented',
  });
});

export default router;
