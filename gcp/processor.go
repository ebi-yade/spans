// package gcp is an abbreviation for Google Cloud 'Processor', not 'Platform'.
package gcp

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Processor is a custom SpanProcessor that ensures that the attributes of a span are compatible with Google Cloud Trace
type Processor struct {
	nextProcessor trace.SpanProcessor
}

// NewProcessor creates a new Processor
func NewProcessor(next trace.SpanProcessor) *Processor {
	return &Processor{
		nextProcessor: next,
	}
}

// OnStart is called when a span starts
func (p *Processor) OnStart(parent context.Context, s trace.ReadWriteSpan) {
	p.convertAttributes(s)
	p.nextProcessor.OnStart(parent, s)
}

// OnEnd is called when a span ends
func (p *Processor) OnEnd(s trace.ReadOnlySpan) {
	p.nextProcessor.OnEnd(s)
}

// Shutdown shuts down the processor.
func (p *Processor) Shutdown(ctx context.Context) error {
	return p.nextProcessor.Shutdown(ctx)
}

// ForceFlush forces the processor to flush any buffered spans.
func (p *Processor) ForceFlush(ctx context.Context) error {
	return p.nextProcessor.ForceFlush(ctx)
}

// convertAttributes ensures that the attributes of a span are compatible with Google Cloud Trace
func (p *Processor) convertAttributes(s trace.ReadWriteSpan) {
	attrs := s.Attributes()
	overwrittenAttrs := make([]attribute.KeyValue, 0, len(attrs))

	for _, attr := range attrs {
		compatible, isDirty := p.ensureCompatibleAttr(attr)
		if isDirty {
			overwrittenAttrs = append(overwrittenAttrs, compatible)
		}
	}

	s.SetAttributes(overwrittenAttrs...)
}

// ensureCompatibleAttr stringifies the attribute value with types that are supported by Google Cloud Trace
func (p *Processor) ensureCompatibleAttr(attr attribute.KeyValue) (attribute.KeyValue, bool) {
	key := attr.Key

	switch attr.Value.Type() {
	case attribute.FLOAT64:
		return key.String(fmt.Sprintf("%f", attr.Value.AsFloat64())), true

	case attribute.BOOLSLICE:
		return key.String(formatSlice(attr.Value.AsBoolSlice(), nakedElement)), true
	case attribute.INT64SLICE:
		return key.String(formatSlice(attr.Value.AsInt64Slice(), nakedElement)), true
	case attribute.FLOAT64SLICE:
		return key.String(formatSlice(attr.Value.AsFloat64Slice(), nakedElement)), true
	case attribute.STRINGSLICE:
		return key.String(formatSlice(attr.Value.AsStringSlice(), quotedElement)), true

	default:
		return attr, false
	}
}

// =================================================================================
// Slice formatter, which might be exported and provided as options in the future
// =================================================================================

type ElementFormatter func(elem any) string

func nakedElement(elem any) string {
	return fmt.Sprintf("%v", elem)
}

func quotedElement(elem any) string {
	return fmt.Sprintf(`"%v"`, elem)
}

// formatSlice is the default implementation to format a slice of elements
func formatSlice[t any](slice []t, formatter ElementFormatter) string {
	strSlice := make([]string, 0, len(slice))
	for _, v := range slice {
		strSlice = append(strSlice, formatter(v))
	}

	return "[" + strings.Join(strSlice, ", ") + "]"
}
