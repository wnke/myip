package myip

import (
	"errors"
	"net/netip"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestNewIPDiscover(t *testing.T) {
	ipd, err := NewIPDiscover()
	assert.Nil(t, err)
	assert.NotNil(t, ipd)

	assert.Equal(t, 0, ipd.index)
	assert.Equal(t, DEFAULT_PROVIDERS, ipd.Providers())
}

func TestNewIPDiscoverWithProvidersErrEmpty(t *testing.T) {
	providers := []string{}
	ipd, err := NewIPDiscoverWithProviders(providers)
	assert.NotNil(t, err)
	assert.Nil(t, ipd)

	assert.ErrorIs(t, err, ErrEmptyProviders)
}

func TestNewIPDiscoverWithProvidersErrInvalidURL(t *testing.T) {
	providers := []string{"batatas1\123"}
	ipd, err := NewIPDiscoverWithProviders(providers)
	assert.NotNil(t, err)
	assert.Nil(t, ipd)

	assert.ErrorIs(t, err, ErrInvalidURL)
}
func TestNextProvider(t *testing.T) {
	providers := append(DEFAULT_PROVIDERS, DEFAULT_PROVIDERS...)

	ipd, err := NewIPDiscover()
	assert.Nil(t, err)
	assert.NotNil(t, ipd)

	for _, provider := range providers {
		assert.Equal(t, provider, ipd.nextProvider())
	}
}
func TestDiscover(t *testing.T) {
	ipd, err := NewIPDiscover()
	assert.Nil(t, err)
	assert.NotNil(t, ipd)

	httpmock.ActivateNonDefault(ipd.client)
	defer httpmock.DeactivateAndReset()

	expectedIP := netip.MustParseAddr("127.0.0.1")
	httpmock.RegisterResponder("GET", ipd.Providers()[ipd.index],
		httpmock.NewStringResponder(200, expectedIP.String()))

	ip, err := ipd.Discover()
	assert.Nil(t, err)
	assert.NotNil(t, ip)

	assert.Equal(t, &expectedIP, ip)
	assert.Equal(t, 1, ipd.index)
}

func TestDiscoverErrNoResponse(t *testing.T) {
	ipd, err := NewIPDiscover()
	assert.Nil(t, err)
	assert.NotNil(t, ipd)

	httpmock.ActivateNonDefault(ipd.client)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", ipd.Providers()[ipd.index],
		httpmock.NewErrorResponder(errors.New("failed")))

	ip, err := ipd.Discover()
	assert.NotNil(t, err)
	assert.Nil(t, ip)

	assert.ErrorIs(t, err, ErrNoResponse)
}

func TestDiscoverErrStatusCode(t *testing.T) {
	ipd, err := NewIPDiscover()
	assert.Nil(t, err)
	assert.NotNil(t, ipd)

	httpmock.ActivateNonDefault(ipd.client)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", ipd.Providers()[ipd.index],
		httpmock.NewStringResponder(403, "Forbidden"))

	ip, err := ipd.Discover()
	assert.NotNil(t, err)
	assert.Nil(t, ip)

	assert.ErrorIs(t, err, ErrStatusCode)
}

func TestDiscoverErrParseResponse(t *testing.T) {
	ipd, err := NewIPDiscover()
	assert.Nil(t, err)
	assert.NotNil(t, ipd)

	httpmock.ActivateNonDefault(ipd.client)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", ipd.Providers()[ipd.index],
		httpmock.NewStringResponder(200, "batatas"))

	ip, err := ipd.Discover()
	assert.NotNil(t, err)
	assert.Nil(t, ip)

	assert.ErrorIs(t, err, ErrParseResponse)
}
