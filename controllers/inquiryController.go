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
	if apiKey == "" {
		c.JSON(500, gin.H{"error": "RESEND_API_KEY not set"})
		return
	}

	client := resend.NewClient(apiKey)

	adminHTML := `
<!DOCTYPE html>
<html>
<body style="margin:0;padding:0;background:#f4f4f4;">
<table width="100%" cellpadding="0" cellspacing="0" bgcolor="#f4f4f4">
<tr>
<td align="center">

<table width="600" cellpadding="20" cellspacing="0" bgcolor="#ffffff" style="font-family:Arial;">
<tr>
<td>

<h2 style="margin-top:0;color:#111;">📥 Inquiry Baru</h2>

<p><strong>Nama:</strong> ` + inquiry.Name + `</p>
<p><strong>Email:</strong> ` + inquiry.Email + `</p>
<p><strong>Contact:</strong> ` + inquiry.Contact + `</p>

<p><strong>Pesan:</strong></p>
<p>` + inquiry.Message + `</p>

<br>

<a href="mailto:` + inquiry.Email + `" 
style="background:#0f0f0f;color:#ffffff;padding:10px 15px;text-decoration:none;">
Balas Email
</a>

</td>
</tr>
</table>

</td>
</tr>
</table>
</body>
</html>
`

	adminParams := &resend.SendEmailRequest{
		From:    "Weavory Studio <noreply@weavorystudio.com>",
		To:      []string{"weavorystudio@gmail.com"},
		Subject: "Inquiry Baru dari Website",
		Html:    adminHTML,
		Text:    "Inquiry baru dari " + inquiry.Name,
	}

	resAdmin, err := client.Emails.Send(adminParams)
	if err != nil {
		println("Admin email failed:", err.Error())
	} else {
		println("Admin email sent:", resAdmin.Id)
	}

	userHTML := `
<!DOCTYPE html>
<html>
<body style="margin:0;padding:0;background:#102F76;">
<table width="100%" cellpadding="0" cellspacing="0" bgcolor="#102F76">
<tr>
<td align="center">

<table width="600" cellpadding="0" cellspacing="0" bgcolor="#ffffff" style="font-family:Arial;">

<!-- HEADER -->
<tr>
<td align="center" bgcolor="#102F76" style="padding:25px;">
<img 
src="https://res.cloudinary.com/dphlgt5hf/image/upload/f_auto,q_auto,w_200/Design_2-removebg-preview_np1r2i.png" 
width="120"
style="display:block;margin:0 auto;" 
/>
<p style="color:#477D7B;font-size:12px;margin-top:10px;letter-spacing:1px;">
PREMIUM TAILOR & FASHION STUDIO
</p>
</td>
</tr>

<!-- CONTENT -->
<tr>
<td style="padding:30px;">

<h2 style="margin:0;color:#102F76;">
Terima kasih, ` + inquiry.Name + `
</h2>

<p style="color:#555;">
Kami telah menerima inquiry Anda:
</p>

<ul style="color:#333;">
<li><b>Nama:</b> ` + inquiry.Name + `</li>
<li><b>Email:</b> ` + inquiry.Email + `</li>
<li><b>Contact:</b> ` + inquiry.Contact + `</li>
</ul>

<p><b>Pesan:</b></p>
<p>` + inquiry.Message + `</p>

<br>

<p style="color:#555;">
Tim kami akan segera menghubungi Anda.
</p>

<div style="text-align:center;margin-top:25px;">
<a href="https://weavorystudio.com"
style="background:#477D7B;color:#ffffff;padding:12px 20px;text-decoration:none;border-radius:4px;">
Kunjungi Website
</a>
</div>

</td>
</tr>

<!-- FOOTER -->
<tr>
<td align="center" bgcolor="#102F76" style="padding:20px;color:#ffffff;font-size:12px;">
Weavory Studio<br/>
<span style="color:#477D7B;">
Custom Apparel • Tailor • Sablon
</span>
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
		Subject: "Terima kasih sudah menghubungi Weavory Studio",
		Html:    userHTML,
		Text: "Halo " + inquiry.Name + ", terima kasih sudah menghubungi Weavory Studio. Kami akan segera menghubungi Anda.",
	}

	resUser, err := client.Emails.Send(userParams)
	if err != nil {
		println("User email failed:", err.Error())
	} else {
		println("User email sent:", resUser.Id)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inquiry sent successfully",
		"data":    inquiry,
	})
}