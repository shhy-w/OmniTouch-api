package mail

// VerificationCodeTemplate
// fmt.Sprintf(VerificationCodeTemplate, title, code, action, validTime)
const VerificationCodeTemplate = `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Verification Code</title>
    <style>
      body {
        font-family: Arial, sans-serif;
        background-color: #f4f4f4;
        margin: 0;
        padding: 0;
        color: #333333;
      }
      .container {
        max-width: 600px;
        margin: 50px auto;
        background-color: #ffffff;
        padding: 16px;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      }
      .header {
        text-align: center;
        padding-bottom: 16px;
        border-bottom: 1px solid #eeeeee;
      }
      .header h1 {
        margin: 0;
        font-size: 24px;
        color: #007bff;
      }
      .content {
        text-align: center;
        padding: 10px 0;
      }
      .content p {
        font-size: 18px;
        line-height: 1.3;
        margin: 4px 0;
      }
      .code {
        display: inline-block;
        padding: 10px 24px;
        font-size: 24px;
        color: #ffffff;
        background-color: #1282f9;
        border-radius: 5px;
        letter-spacing: 2px;
        margin: 16px 0;
      }
      .footer {
        text-align: center;
        font-size: 14px;
        color: #777777;
      }
      .footer p {
        margin: 12px 0;
      }
      .footer a {
        color: #007bff;
        text-decoration: none;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="header">
        <h1>%s</h1>
      </div>
      <div class="content">
        <p>Hello,</p>
        <p>Your verification code is:</p>
        <div class="code">%s</div>
        <p>
          Please use this code to complete your %s.
          <br />
          The code is valid for <strong>%v minute</strong>s.
        </p>
      </div>
      <div class="footer">
        <p>
          If you didn't request this code, please ignore this email or
          <a href="#">contact support</a>.
        </p>
        <p>
          Thank you,<br />
          <span style="display: inline-block; margin-top: 2px">Your <a href="%v">%v</a></span>
        </p>
      </div>
    </div>
  </body>
</html>
`

const RegisterTemplate = `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Verify Your Email</title>
    <style>
      body {
        font-family: Arial, sans-serif;
        background-color: #f4f4f4;
        margin: 0;
        padding: 0;
        color: #333333;
      }
      .container {
        max-width: 600px;
        margin: 50px auto;
        background-color: #ffffff;
        padding: 20px;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      }
      .header {
        text-align: center;
        padding-bottom: 20px;
        border-bottom: 1px solid #eeeeee;
      }
      .header h1 {
        margin: 0;
        font-size: 24px;
        color: #007bff;
      }
      .content {
        max-width: 380px;
        margin: 0 auto;
        text-align: center;
        padding: 30px 0;
        line-height: 1.3;
      }
      .content p {
        font-size: 18px;
        margin: 20px 0;
      }
      a {
        color: #007bff;
      }
      .button {
        display: inline-block;
        padding: 12px 30px;
        font-size: 16px;
        color: #ffffff !important;
        background-color: #007bff;
        border-radius: 5px;
        text-decoration: none;
      }
      .content .tips {
        font-size: 14px;
      }
      .footer {
        text-align: center;
        font-size: 14px;
        color: #777777;
        margin-top: 30px;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="header">
        <h1>Verify your email address</h1>
      </div>
      <div class="content">
        <p>Hello,</p>
        <p>
          Thank you for registering. Please click the button below to verify
          your email address:
        </p>
        <a href="%s" class="button">Verify email</a>
        <p class="tips">
          If you did not register, please ignore this email. The link will be
          expired after <strong>%v minutes</strong>
        </p>
      </div>
      <div class="footer">
        <p>Thank you,<br />Your <a href="%v">%v</a></p>
      </div>
    </div>
  </body>
</html>
`

const ResetPasswordTemplate = `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Reset Password Successfully</title>
    <style>
      body {
        font-family: Arial, sans-serif;
        background-color: #f4f4f4;
        margin: 0;
        padding: 0;
        color: #333333;
      }
      .container {
        max-width: 600px;
        margin: 50px auto;
        background-color: #ffffff;
        padding: 20px;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      }
      .header {
        text-align: center;
        padding-bottom: 20px;
        border-bottom: 1px solid #eeeeee;
      }
      .header h1 {
        margin: 0;
        font-size: 24px;
        color: #007bff;
      }
      .content {
        max-width: 380px;
        margin: 0 auto;
        text-align: center;
        padding: 30px 0;
        line-height: 1.3;
      }
      .content p {
        font-size: 18px;
        margin: 20px 0;
      }
      a {
        color: #007bff;
      }
      .content .tips {
        font-size: 14px;
      }
      .footer {
        text-align: center;
        font-size: 14px;
        color: #777777;
        margin-top: 30px;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="header">
        <h1>Reset password successfully</h1>
      </div>
      <div class="content">
        <p>Hello,</p>
        <p>
          Your password has been reset.
        </p>
        <p class="tips">
          If you did not request this change, please reset your password immediately.
        </p>
      </div>
      <div class="footer">
        <p>Thank you,<br />Your <a href="%v">%v</a></p>
      </div>
    </div>
  </body>
</html>
`
