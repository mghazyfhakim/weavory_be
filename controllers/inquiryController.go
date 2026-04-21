package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/resend/resend-go/v2"

	"weavory-backend/config"
	"weavory-backend/models"
	"weavory-backend/utils"
)

func CreateInquiry(c *gin.Context) {
	var inquiry models.Inquiry

	if err := c.BindJSON(&inquiry); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if inquiry.Name == "" || inquiry.Email == "" || inquiry.Contact == "" {
		utils.Error(c, 400, "Semua field wajib diisi")
		return
	}

	_, err := config.DB.Exec(
		"INSERT INTO inquiries (name,email,contact,message) VALUES ($1,$2,$3,$4)",
		inquiry.Name,
		inquiry.Email,
		inquiry.Contact,
		inquiry.Message,
	)

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	apiKey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(apiKey)

	htmlContent := `
<!DOCTYPE html>
<html>
<body style="margin:0;background:#f4f4f4;font-family:Arial;">

<table width="100%" cellpadding="20">
<tr>
<td align="center">

<table width="600" style="background:#fff;border-radius:10px;padding:20px;">

<tr>
<td>
<h2 style="margin-top:0;">📥 Inquiry Baru</h2>
<p style="color:#666;">Seseorang mengirim pesan dari website:</p>

<table width="100%" style="margin-top:15px;font-size:14px;">
<tr><td><strong>Nama</strong></td><td>` + inquiry.Name + `</td></tr>
<tr><td><strong>Email</strong></td><td>` + inquiry.Email + `</td></tr>
<tr><td><strong>Contact</strong></td><td>` + inquiry.Contact + `</td></tr>
</table>

<div style="margin-top:15px;padding:15px;background:#fafafa;border-radius:8px;">
<strong>Pesan:</strong>
<p style="margin:5px 0;">` + inquiry.Message + `</p>
</div>

<div style="margin-top:20px;">
<a href="mailto:` + inquiry.Email + `" 
style="background:#111;color:#fff;padding:10px 15px;text-decoration:none;border-radius:5px;">
Balas Sekarang
</a>
</div>

</td>
</tr>

</table>

</td>
</tr>
</table>

</body>
</html>
`

	params := &resend.SendEmailRequest{
		From:    "Weavory Studio <inquiry@weavorystudio.com>",
		To:      []string{"weavorystudio@gmail.com"},
		Subject: "New Inquiry from Website",
		Html:    htmlContent,
	}

	_, err = client.Emails.Send(params)
	if err != nil {
		println("Email failed:", err.Error())
	}

	userHtml := `
<!DOCTYPE html>
<html>
<body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Helvetica,Arial,sans-serif;">

<table width="100%" cellpadding="0" cellspacing="0" style="padding:30px 0;background:#f4f4f4;">
  <tr>
    <td align="center">

      <!-- CONTAINER -->
      <table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 4px 20px rgba(0,0,0,0.08);">

        <!-- HEADER -->
        <tr>
          <td style="background:#0f0f0f;padding:25px;text-align:center;">
            <img src="LOGO_URL" alt="Weavory Studio" style="height:40px;margin-bottom:10px;">
            <p style="color:#c9a96e;font-size:12px;letter-spacing:2px;margin:0;">
              PREMIUM TAILOR & FASHION STUDIO
            </p>
          </td>
        </tr>

        <!-- HERO TEXT -->
        <tr>
          <td style="padding:35px 40px 10px 40px;text-align:center;">
            <h1 style="margin:0;font-size:24px;color:#111;">
              Terima kasih, ` + inquiry.Name + `
            </h1>
            <p style="color:#777;margin-top:10px;font-size:14px;">
              Inquiry Anda telah kami terima dengan baik
            </p>
          </td>
        </tr>

        <!-- DIVIDER -->
        <tr>
          <td style="padding:10px 40px;">
            <hr style="border:none;border-top:1px solid #eee;">
          </td>
        </tr>

        <!-- CONTENT -->
        <tr>
          <td style="padding:10px 40px 30px 40px;">
            
            <p style="font-size:14px;color:#555;">
              Berikut detail yang Anda kirimkan:
            </p>

            <table width="100%" style="margin-top:15px;border-collapse:collapse;font-size:14px;">
              <tr>
                <td style="padding:10px 0;color:#999;">Nama</td>
                <td style="padding:10px 0;color:#111;text-align:right;">` + inquiry.Name + `</td>
              </tr>
              <tr>
                <td style="padding:10px 0;color:#999;">Email</td>
                <td style="padding:10px 0;color:#111;text-align:right;">` + inquiry.Email + `</td>
              </tr>
              <tr>
                <td style="padding:10px 0;color:#999;">Contact</td>
                <td style="padding:10px 0;color:#111;text-align:right;">` + inquiry.Contact + `</td>
              </tr>
            </table>

            <!-- MESSAGE BOX -->
            <div style="margin-top:20px;padding:15px;background:#fafafa;border-radius:8px;border:1px solid #eee;">
              <p style="margin:0;color:#777;font-size:13px;">Pesan Anda:</p>
              <p style="margin:5px 0 0;color:#111;font-size:14px;">
                ` + inquiry.Message + `
              </p>
            </div>

            <p style="margin-top:25px;font-size:14px;color:#555;">
              Tim kami akan segera menghubungi Anda dalam waktu dekat.
            </p>

            <!-- CTA -->
            <div style="text-align:center;margin-top:30px;">
              <a href="https://weavorystudio.com"
                 style="background:#c9a96e;color:#111;padding:12px 24px;text-decoration:none;border-radius:6px;font-weight:bold;">
                EXPLORE OUR COLLECTION
              </a>
            </div>

          </td>
        </tr>

        <!-- FOOTER -->
        <tr>
          <td style="background:#0f0f0f;color:#aaa;text-align:center;padding:25px;font-size:12px;">
            <p style="margin:0;color:#fff;font-weight:bold;">Weavory Studio</p>
            <p style="margin:5px 0;">Custom Apparel • Tailor • Sablon</p>

            <p style="margin-top:10px;">
              <a href="https://weavorystudio.com" style="color:#c9a96e;text-decoration:none;">Website</a> •
              <a href="#" style="color:#c9a96e;text-decoration:none;">Instagram</a>
            </p>

            <p style="margin-top:15px;color:#666;">
              © 2026 Weavory Studio. All rights reserved.
            </p>
          </td>
        </tr>

      </table>

    </td>
  </tr>
</table>

</body>
</html>
`

	userParams := &resend.SendEmailRequest{
		From:    "Weavory Studio <noreply@weavorystudio.com>",
		To:      []string{inquiry.Email},
		Subject: "Konfirmasi Inquiry - Weavory Studio",
		Html:    userHtml,
	}

	go func() {
		_, err := client.Emails.Send(userParams)
		if err != nil {
			println("Auto-reply failed:", err.Error())
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "Inquiry sent successfully",
		"data":    inquiry,
	})
}
