import { Router } from 'express';

const router = Router();

/**
 * @openapi
 * /api/bills:
 *   get:
 *     tags:
 *       - Bills
 *     summary: Get list of bills
 *     parameters:
 *       - in: query
 *         name: page
 *         schema:
 *           type: integer
 *       - in: query
 *         name: limit
 *         schema:
 *           type: integer
 *       - in: query
 *         name: status
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: List of bills
 */
router.get('/', (_req, res) => {
  res.status(200).json({
    success: true,
    data: [],
    message: 'Get bills endpoint - to be implemented',
  });
});

export default router;
