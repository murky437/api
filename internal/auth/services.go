package auth

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

func isPasswordValid(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

type GuestLoginLimiter struct {
	emailLimiters sync.Map
	ipLimiters    sync.Map
}

func NewGuestLoginLimiter() *GuestLoginLimiter {
	return &GuestLoginLimiter{}
}

func (l *GuestLoginLimiter) getEmailLimiter(email string) *rate.Limiter {
	lim, _ := l.emailLimiters.LoadOrStore(email, rate.NewLimiter(rate.Every(60*time.Second), 1))
	return lim.(*rate.Limiter)
}

func (l *GuestLoginLimiter) getIpLimiter(ip string) *rate.Limiter {
	lim, _ := l.ipLimiters.LoadOrStore(ip, rate.NewLimiter(rate.Every(30*time.Second), 1))
	return lim.(*rate.Limiter)
}

func (l *GuestLoginLimiter) getIp(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ip
}

func (l *GuestLoginLimiter) Reserve(email string, r *http.Request) (allowed bool, retryAfterSeconds int) {
	emailLimiter := l.getEmailLimiter(email)
	ipLimiter := l.getIpLimiter(l.getIp(r))

	if emailLimiter.Tokens() >= 1 && ipLimiter.Tokens() >= 1 {
		emailLimiter.Allow()
		ipLimiter.Allow()
		return true, 0
	}

	emailRetrySeconds := 0
	if emailLimiter.Tokens() < 1 {
		emailRetrySeconds = int((1-emailLimiter.Tokens())/float64(emailLimiter.Limit()) + 0.5)
	}

	ipRetrySeconds := 0
	if ipLimiter.Tokens() < 1 {
		ipRetrySeconds = int((1-ipLimiter.Tokens())/float64(ipLimiter.Limit()) + 0.5)
	}

	if ipRetrySeconds > emailRetrySeconds {
		return false, ipRetrySeconds
	}
	return false, emailRetrySeconds
}
