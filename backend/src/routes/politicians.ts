import { Router } from 'express';

const router = Router();

/**
 * @openapi
 * /api/politicians:
 *   get:
 *     tags:
 *       - Politicians
 *     summary: Get list of politicians
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
 *         name: party
 *         schema:
 *           type: string
 *       - in: query
 *         name: state
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: List of politicians
 */
router.get('/', (_req, res) => {
  res.status(200).json({
    success: true,
    data: [],
    message: 'Get politicians endpoint - to be implemented',
  });
});

/**
 * @openapi
 * /api/politicians/{id}:
 *   get:
 *     tags:
 *       - Politicians
 *     summary: Get politician by ID
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Politician details
 *       404:
 *         description: Politician not found
 */
router.get('/:id', (_req, res) => {
  res.status(200).json({
    success: true,
    message: 'Get politician by ID endpoint - to be implemented',
  });
});

export default router;
