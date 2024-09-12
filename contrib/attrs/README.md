## Example

```go
package main

import (
	"context"

	"github.com/mashiike/otelattr/contrib/attrs"
)

type Person struct {
	ID      string `otelattr:"id"`
	Name    string `otelattr:"name"`
	Age     int    `otelattr:"age"`
	Address string `otelattr:"-"`
}

func GetPeopleFromTeam(ctx context.Context, teamID string) ([]Person, error) {
	ctx, span := tracer.Start(ctx, "GetPeopleFromTeam", attrs.OnStart(attrs.String("team.id", teamID)))
	defer span.End()

	people, err := fetchPeople(ctx, teamID)
	if err != nil {
		return nil, errors.Wrap(err, "error fetchPeople")
	}
	for _, p := range people {
		attrs.Set(span, attrs.Any("person", p))
	}

	return people, nil
}

```
