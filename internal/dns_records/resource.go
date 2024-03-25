// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_records

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/cloudflare/cloudflare-go"
	"github.com/cloudflare/cloudflare-go/dns"
	"github.com/cloudflare/cloudflare-go/option"
	"github.com/cloudflare/cloudflare-terraform/internal/apijson"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DNSRecordsResource{}

func NewResource() resource.Resource {
	return &DNSRecordsResource{}
}

// DNSRecordsResource defines the resource implementation.
type DNSRecordsResource struct {
	client *cloudflare.Client
}

func (r *DNSRecordsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_records"
}

func (r *DNSRecordsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cloudflare.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"unexpected resource configure type",
			fmt.Sprintf("Expected *cloudflare.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DNSRecordsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DNSRecordsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	dataBytes, err := apijson.Marshal(data)
	if err != nil {
		resp.Diagnostics.AddError("failed to create resource", err.Error())
		return
	}
	res := new(http.Response)
	env := DNSRecordsResultEnvelope{*data}
	_, err = r.client.DNS.Records.New(
		ctx,
		dns.RecordNewParamsDNSRecordsARecord{
			ZoneID: cloudflare.F(data.ZoneID.ValueString()),
		},
		option.WithRequestBody("application/json", dataBytes),
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(loggingMiddleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	bytes, _ := io.ReadAll(res.Body)
	tflog.Warn(ctx, fmt.Sprintf("BYTES %s", string(bytes)))
	apijson.Unmarshal(bytes, &env)
	data = &env.Result

	tflog.Warn(ctx, fmt.Sprintf("ENV %#v", env))
	tflog.Warn(ctx, fmt.Sprintf("DATA %#v", data))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSRecordsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DNSRecordsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	env := DNSRecordsResultEnvelope{*data}
	_, err := r.client.DNS.Records.Get(
		ctx,
		data.ID.ValueString(),
		dns.RecordGetParams{
			ZoneID: cloudflare.F(data.ZoneID.ValueString()),
		},
		option.WithResponseBodyInto(&env),
		option.WithMiddleware(loggingMiddleware(ctx)),
	)
	data = &env.Result

	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSRecordsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DNSRecordsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	dataBytes, err := apijson.Marshal(data)
	if err != nil {
		resp.Diagnostics.AddError("failed to create resource", err.Error())
		return
	}
	res := new(http.Response)
	env := DNSRecordsResultEnvelope{*data}
	_, err = r.client.DNS.Records.Update(
		ctx,
		data.ID.ValueString(),
		dns.RecordUpdateParamsDNSRecordsARecord{
			ZoneID: cloudflare.F(data.ZoneID.ValueString()),
		},
		option.WithRequestBody("application/json", dataBytes),
		option.WithResponseBodyInto(&res),
		option.WithMiddleware(loggingMiddleware(ctx)),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	bytes, _ := io.ReadAll(res.Body)
	tflog.Warn(ctx, fmt.Sprintf("BYTES %s", string(bytes)))
	apijson.Unmarshal(bytes, &env)
	data = &env.Result

	tflog.Warn(ctx, fmt.Sprintf("ENV %#v", env))
	tflog.Warn(ctx, fmt.Sprintf("DATA %#v", data))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSRecordsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DNSRecordsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DNS.Records.Delete(
		ctx,
		data.DNSRecordID.ValueString(),
		dns.RecordDeleteParams{
			ZoneID: cloudflare.F(data.ZoneID.ValueString()),
		},
		option.WithMiddleware(loggingMiddleware(ctx)),
	)

	if err != nil {
		resp.Diagnostics.AddError("failed to make http request", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func loggingMiddleware(ctx context.Context) option.Middleware {
	return func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		logRequest(ctx, req)

		resp, err := next(req)

		logResponse(ctx, resp)

		return resp, err
	}
}

func logRequest(ctx context.Context, req *http.Request) error {
	// Log headers
	tflog.Warn(ctx, "Headers:")
	for name, values := range req.Header {
		for _, value := range values {
			tflog.Warn(ctx, fmt.Sprintf("%s: %s\n", name, value))
		}
	}

	if req.Body != nil {
		// Read the body without mutating the original response
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}

		// Restore the original body to the response so it can be read again
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Log the body
		tflog.Warn(ctx, fmt.Sprintf("Body: %s\n", string(bodyBytes)))
	}

	return nil
}

func logResponse(ctx context.Context, resp *http.Response) error {
	// Log the status code
	tflog.Warn(ctx, fmt.Sprintf("Status: %s\n", resp.Status))

	// Log headers
	tflog.Warn(ctx, "Headers:")
	for name, values := range resp.Header {
		for _, value := range values {
			tflog.Warn(ctx, fmt.Sprintf("%s: %s\n", name, value))
		}
	}

	// Read the body without mutating the original response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Restore the original body to the response so it can be read again
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Log the body
	tflog.Warn(ctx, fmt.Sprintf("Body: %s\n", string(bodyBytes)))

	return nil
}
