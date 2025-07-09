package jwt_key_mock

import "errors"

func GetPrivateKey() string {
	return `-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCX6qX6hrE6LviN
nbtgwVL0Mxsz1EHlYsTGvnQka3BiwfgnzviNeV3iSO0m3j0jiF7Hx96l+BVl0m/0
iYxidFwtkn2oLh7FmvdWQ8Z1o91V1FkS5hF5LW9WeSy3yVdnyWnZNCbX9WaG5DOX
iiJPEZdGetSgHwK1Sg05m9gVqw3Tuj2Cv4aN5h9KR7h2BLR/XMvuXBztrHSamKL4
b/kaUCfCewlH7kFtvA/EY1HhiKN40ZsvxHz07ZH4Z35IwMNarqhDsX6jHPMXj9yN
zDRTG9Vbo0c0YrRRZSTWz642jdIwwah8BFqjDRBI4t6gDyykoubN6xFuwl2lVOix
0b7wkDvfAgMBAAECggEALwGRqxO7SfQXv4YTYSxglSQtIhvm6enmVlK/EDfYVg6L
0RGHMgtgQ/DpxJGWnGSJY0rSd0tsn41+S821zQ8RKJ6/1+axadZ5+szM30C9b13d
0+2xcWvgHAMxlYqsy3X0lRtAe6k4uLtqxVSeKhYY11uRaBxAU8Uek3lH4uiDRkRL
pTWJ1EmiTUATmN+nvFDYD4e6ujqBms++fyRYeAcbGgsv+/6PpbSb3yCjLhvUXz7i
k9LbKzB8rXVzus37i+O6uRKve/kLGuZ6zfxNpoDmnuxKItHFqN9QExWFnecHj9sj
fMavFXpE/FEoudcY9lrg5IN+pNnVkQGDcLWwGWXuZQKBgQDQHCeerUE5F+AbRrRN
u9H96Ft+38SIXNMxO9zD/hnTpOlO9uT/PByZq7AdYzrKSDqE2yerDpLPHuulHgv+
OqgE3Csu5Zod7WhfzdT/DYWOiiog3NfotNsezpX5YdLWdQSlZ2qOe8mHBA7cLVW8
jagf2e7nMVjy+M9qihP5RLQUgwKBgQC64B41Om2PUyqFZSczcMOnNKHwW75t3NaT
581kGvyaTeDhmC+FUDogNgUGZJngwDGQyfyrWEx1pzgHV8pymi47kqGUPLsuCjBN
f3X9xZ1h6xhfMZaxpeGkfhPWZ/B/b9Uu3L/F2hMvmrV1XCkrv8yAAARujsgJQt/w
ilWxvJL0dQKBgDHzMAdD6m27r0ycsdYeiI564MsZBmD8dqxQg/J+4NANuvn7BIfG
c87miITlNk0q/PC6cVD7VH1mHIUrKxHAHmfcOHkvHsikHPMxwjfdlPrbarUsjJ4M
GrPQPer3cdWLjKvuoILGb156uN5b+0IgdgP/GPpgu8rFsXMO5TBlLxvlAoGALSWi
BqgD+gFUn3+Nle7jRc0AZoozmmUk7fytcUbXyguQjc/vgxybvlZuplm9lz+3ecxi
n56ocjAg6B08iq1XCAtnv+FgM0JA4ygtAE8ys4pRjAX16xsxRUU0U7Mutgr1jOnF
5u3FftW4iw7l32zp4e6fI3qZNyuR4JH7HAJ72lUCgYAQaYKvzEj2dPLjBiMiP9TS
Yyqpa5zKXazua1EmgQuvK6pMALXmHNGc3sGEOv+81tRyxiBgsrdTK7aWIk8KpgBV
ogXbEuCspmBr/C5CGy+W8t0XCPNTgFUH9i6VwKcHS37vePKUOiaHK26GlbjJdfZ7
5gl8yTGteIufRn6q4Nj6Jw==
-----END PRIVATE KEY-----`
}

func GetPublicKey() string {
	return `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAl+ql+oaxOi74jZ27YMFS
9DMbM9RB5WLExr50JGtwYsH4J874jXld4kjtJt49I4hex8fepfgVZdJv9ImMYnRc
LZJ9qC4exZr3VkPGdaPdVdRZEuYReS1vVnkst8lXZ8lp2TQm1/VmhuQzl4oiTxGX
RnrUoB8CtUoNOZvYFasN07o9gr+GjeYfSke4dgS0f1zL7lwc7ax0mpii+G/5GlAn
wnsJR+5BbbwPxGNR4YijeNGbL8R89O2R+Gd+SMDDWq6oQ7F+oxzzF4/cjcw0UxvV
W6NHNGK0UWUk1s+uNo3SMMGofARaow0QSOLeoA8spKLmzesRbsJdpVTosdG+8JA7
3wIDAQAB
-----END PUBLIC KEY-----`
}

// MockSigner forces an error when signing JWTs.
type MockSigner struct{}

func (m *MockSigner) Sign(signingString string, key interface{}) ([]byte, error) {
	return []byte{}, errors.New("forced signing error")
}

func (m *MockSigner) Alg() string {
	return "RS256"
}

func (m *MockSigner) Verify(signingString string, sig []byte, key interface{}) error {
	return nil
}
