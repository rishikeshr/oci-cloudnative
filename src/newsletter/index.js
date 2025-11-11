/**
 * Copyright © 2019, Oracle and/or its affiliates. All rights reserved.
 * The Universal Permissive License (UPL), Version 1.0
 */
// EDOT instrumentation must be loaded first
require('./tracing');

const express = require('express');
const bodyParser = require('body-parser');
const helmet = require('helmet');
const morgan = require('morgan');
const nodemailer = require('nodemailer');

const app = express();
const PORT = process.env.PORT || 3000;

// Middleware
app.use(helmet());
app.use(morgan('combined'));
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

// SMTP Configuration from environment variables
const {
  SMTP_USER,
  SMTP_PASSWORD,
  SMTP_HOST = 'localhost',
  SMTP_PORT = '1025',
  APPROVED_SENDER_EMAIL = 'noreply@mushop.com',
} = process.env;

// Validate required environment variables
const validateConfig = () => {
  if (!SMTP_USER || !SMTP_PASSWORD) {
    console.warn('Warning: SMTP_USER or SMTP_PASSWORD not set. Email functionality may not work.');
  }
  console.log(`Newsletter service configured with SMTP: ${SMTP_HOST}:${SMTP_PORT}`);
};

validateConfig();

// Create nodemailer transporter
const createTransporter = () => {
  return nodemailer.createTransport({
    host: SMTP_HOST,
    port: parseInt(SMTP_PORT, 10),
    secure: false, // true for 465, false for other ports
    auth: SMTP_USER && SMTP_PASSWORD ? {
      user: SMTP_USER,
      pass: SMTP_PASSWORD,
    } : undefined,
    tls: {
      rejectUnauthorized: false, // Allow self-signed certificates
    },
  });
};

// Send subscription email
const sendSubscriptionEmail = async (toEmail) => {
  const transporter = createTransporter();

  const mailOptions = {
    from: APPROVED_SENDER_EMAIL,
    to: toEmail,
    subject: 'Welcome to MuShop Newsletter! 🎉',
    html: `
      <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
        <h1 style="color: #4CAF50;">Thanks for Subscribing!</h1>
        <p>Hello,</p>
        <p>Thank you for confirming your <strong>subscription</strong> to the MuShop newsletter!</p>
        <p>You'll now receive updates about:</p>
        <ul>
          <li>New product arrivals</li>
          <li>Exclusive deals and discounts</li>
          <li>Special promotions</li>
        </ul>
        <p>Stay tuned for exciting updates!</p>
        <hr style="border: 1px solid #eee; margin: 20px 0;">
        <p style="color: #888; font-size: 12px;">
          This is an automated message from MuShop. If you received this in error, please ignore it.
        </p>
      </div>
    `,
    text: 'Thanks for confirming your subscription to the MuShop newsletter!',
  };

  return transporter.sendMail(mailOptions);
};

// Routes

// Health check endpoint
app.get('/health', (req, res) => {
  res.status(200).json({
    status: 'healthy',
    service: 'newsletter',
    timestamp: new Date().toISOString(),
  });
});

// Subscribe endpoint
app.post('/subscribe', async (req, res) => {
  const { email } = req.body;

  // Validate email
  if (!email) {
    return res.status(400).json({
      success: false,
      error: 'Email address is required',
    });
  }

  // Basic email validation
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRegex.test(email)) {
    return res.status(400).json({
      success: false,
      error: 'Invalid email address format',
    });
  }

  try {
    const info = await sendSubscriptionEmail(email);
    console.log(`Subscription email sent to ${email}, messageId: ${info.messageId}`);

    res.status(200).json({
      success: true,
      message: 'Subscription successful! Check your email.',
      messageId: info.messageId,
    });
  } catch (error) {
    console.error(`Failed to send email to ${email}:`, error);
    res.status(500).json({
      success: false,
      error: 'Failed to send subscription email',
      details: error.message,
    });
  }
});

// Backward compatibility with OCI Functions format
app.post('/newsletter/subscribe', async (req, res) => {
  // Forward to main subscribe endpoint
  req.url = '/subscribe';
  app.handle(req, res);
});

// 404 handler
app.use((req, res) => {
  res.status(404).json({
    error: 'Not found',
    path: req.path,
  });
});

// Error handler
app.use((err, req, res, next) => {
  console.error('Error:', err);
  res.status(500).json({
    error: 'Internal server error',
    message: err.message,
  });
});

// Start server
app.listen(PORT, () => {
  console.log(`Newsletter service listening on port ${PORT}`);
  console.log(`Health check: http://localhost:${PORT}/health`);
  console.log(`Subscribe endpoint: POST http://localhost:${PORT}/subscribe`);
});

// Graceful shutdown
process.on('SIGTERM', () => {
  console.log('SIGTERM received, shutting down gracefully...');
  process.exit(0);
});

process.on('SIGINT', () => {
  console.log('SIGINT received, shutting down gracefully...');
  process.exit(0);
});
