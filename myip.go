package myip

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"time"

	"strings"
)

var DEFAULT_PROVIDERS = []string{
	"https://checkip.amazonaws.com",
	"https://ipconfig.io",
	"https://icanhazip.com",
	"https://ifconfig.me",
}

var (
	ErrEmptyProviders = errors.New("provider list is empty")
	ErrInvalidURL     = errors.New("invalid URL")

	ErrNoResponse      = errors.New("failed to get response from provider")
	ErrStatusCode      = errors.New("provider returned a non OK status code")
	ErrInvalidResponse = errors.New("failed to read IP from provider")
	ErrParseResponse   = errors.New("failed to parse IP from provider")
)

type IPDiscover struct {
	providers []string
	index     int
	client    *http.Client
}

func NewIPDiscover() (*IPDiscover, error) {
	return NewIPDiscoverWithProviders(DEFAULT_PROVIDERS)
}

func NewIPDiscoverWithProviders(providers []string) (*IPDiscover, error) {

	if len(providers) == 0 {
		return nil, fmt.Errorf("%w", ErrEmptyProviders)

	}

	for _, provider := range providers {
		if _, err := url.ParseRequestURI(provider); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidURL, err)
		}
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		var dailer net.Dialer
		return dailer.DialContext(ctx, "tcp4", addr)
	}
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	return &IPDiscover{providers: providers, index: 0, client: client}, nil
}

func (ipd *IPDiscover) Providers() []string {
	providers := []string{}
	for _, provider := range ipd.providers {
		providers = append(providers, strings.Clone(provider))
	}

	return providers
}

func (ipd *IPDiscover) nextProvider() string {
	provider := ipd.providers[ipd.index]
	ipd.index = (ipd.index + 1) % len(ipd.providers)

	return provider
}

func (ipd *IPDiscover) Discover() (*netip.Addr, error) {

	provider := ipd.nextProvider()
	response, err := ipd.client.Get(provider)
	if err != nil {
		return nil, fmt.Errorf("%w (%s): %s", ErrNoResponse, provider, err)
	}
	defer response.Body.Close()

	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOK {
		return nil, fmt.Errorf("%w (%s): %d", ErrStatusCode, provider, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("%w (%s): %s", ErrInvalidResponse, provider, err)
	}

	ip, err := netip.ParseAddr(strings.TrimSpace(string(body)))
	if err != nil {
		return nil, fmt.Errorf("%w (%s): %s", ErrParseResponse, provider, err)
	}

	return &ip, nil
}
