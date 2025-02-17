const express = require('express');
const fs = require('fs');
const router = express.Router();

// Endpoint to get policy
router.get('/', (req, res) => {
    fs.readFile('/etc/cp/conf/local_policy.yaml', 'utf8', (err, data) => {
        if (err) {
            return res.status(500).send('Error reading policy file');
        }
        res.send(data);
    });
});

// Endpoint to save policy
router.post('/', (req, res) => {
    const { policy } = req.body;
    fs.writeFile('/etc/cp/conf/local_policy.yaml', policy, 'utf8', (err) => {
        if (err) {
            return res.status(500).send('Error saving policy file');
        }
        res.send('Policy saved successfully');
    });
});

module.exports = router;
