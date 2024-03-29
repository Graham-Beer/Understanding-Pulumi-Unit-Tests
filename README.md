# Pulumi: Unit Tests in Pulumi with Go

## Introduction
Contains a simple Pulumi project to build an EC2 and unit test written with Go to test two areas of the build **before** Pulumi Up is run.
The two tests contained in the **main_test.go** are:
* Test if the service has tags and a name tag.
* Test if port 22 for ssh is exposed.

## Test results

A run of **Go test -v** command when these elements are exposed or not contained, has the standard output:

```
=== RUN   TestInfrastructure
    main_test.go:56: 
        	Error Trace:	/mnt/c/Projects/Pulumi_ec2_test/main_test.go:56
        	            				/mnt/c/Projects/Pulumi_ec2_test/value.go:586
        	            				/mnt/c/Projects/Pulumi_ec2_test/value.go:370
        	            				/mnt/c/Projects/Pulumi_ec2_test/types.go:494
        	            				/mnt/c/Projects/Pulumi_ec2_test/types.go:606
        	            				/mnt/c/Projects/Pulumi_ec2_test/asm_amd64.s:1598
        	Error:      	Should be false
        	Test:       	TestInfrastructure
        	Messages:   	illegal SSH port 22 open to the Internet (CIDR 0.0.0.0/0) on group urn:pulumi:stack::project::aws:ec2/securityGroup:SecurityGroup::web-secgrp
    main_test.go:37: 
        	Error Trace:	/mnt/c/Projects/Pulumi_ec2_test/main_test.go:37
        	            				/mnt/c/Projects/Pulumi_ec2_test/value.go:586
        	            				/mnt/c/Projects/Pulumi_ec2_test/value.go:370
        	            				/mnt/c/Projects/Pulumi_ec2_test/types.go:494
        	            				/mnt/c/Projects/Pulumi_ec2_test/types.go:606
        	            				/mnt/c/Projects/Pulumi_ec2_test/asm_amd64.s:1598
        	Error:      	map[string]string(nil) does not contain "Name"
        	Test:       	TestInfrastructure
        	Messages:   	missing a Name tag on server urn:pulumi:stack::project::aws:ec2/instance:Instance::web-server-www
--- FAIL: TestInfrastructure (0.01s)
FAIL
exit status 1
FAIL	Pulumi_ec2_test	1.325s
```

Corrections made to the file returns:
```
=== RUN   TestInfrastructure
--- PASS: TestInfrastructure (0.00s)
PASS
ok  	Pulumi_ec2_test	0.977s

```

## Code breakdown

The Pulumi function **WithMocks** has a signature of:  
```go 
func WithMocks(project, stack string, mocks MockResourceMonitor) RunOption
```

The third parameter, **mocks** takes a type of `MockResourceMonitor`. 
`MockResourceMonitor` is of type **Interface**.

Type Interface:
```go
type MockResourceMonitor interface {
	Call(args MockCallArgs) (resource.PropertyMap, error)
	NewResource(args MockResourceArgs) (string, resource.PropertyMap, error)
}
``` 

A new type `mock` is created with two methods of `call` and `NewResource`.
```go
func (mocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return args.Name + "_id", args.Inputs, nil
}

func (mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return args.Args, nil
}
```

The **WithMocks** function is run at the end of the `pulumi.RunErr` function, before closed. 
## Applying the Unit Tests

The opening part of our test function

```go
func TestInfrastructure(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		infra, err := createInfrastructure(ctx)
		assert.NoError(t, err)

		var wg sync.WaitGroup
		wg.Add(2)
```

Package `github.com/stretchr/testify/assert"` provides a set of comprehensive testing tools for use with the normal Go testing system.  
After the the function signature, `func TestInfrastructure(t *testing.T)` the pulumi function, `RunErr`, is used. `RunErr` has a function signature of `func RunErr(body RunFunc, opts ...RunOption) error`.
The `createInfrastructure` function is captured and will be use to test various aspects. 

The `assert.NoError(t, err)` is called to asserts that a function returned no error (i.e. `nil`).

A `WaitGroup` is created from the `sync` package. A WaitGroup waits for a collection of goroutines to finish. The main goroutine calls `Add` to set the number of goroutines to wait for (2 in this example). Then each of the goroutines runs and calls Done when finished. 

> Note - You need to ensure that you call .Add(2) before you execute your goroutine.

### Test 1

```go
pulumi.All(infra.server.URN(), infra.server.Tags).ApplyT(func(all []any) error {
    urn := all[0].(pulumi.URN)
    tags := all[1].(map[string]string)

    assert.Containsf(t, tags, "Name", "missing a Name tag on server %v", urn)
    wg.Done()
    return nil
})
```
If you have multiple outputs and need to join them, the all (`pulumi.All`) function acts like an apply over many resources. The `All` function takes variadic parameter of type `Any`, meaning we can pass any type.  

To access the raw value of an output and transform that value into a new value, use Apply (`ApplyT`). This method accepts a callback that will be invoked with the raw value, once that value is available.

Using the `Containsf` function from the **assert** package, it allows us to:

> Containsf asserts that the specified string, list(array, slice...) or map contains the specified substring or element.

The functions signature is:  
```go
func Containsf(t TestingT, s interface{}, contains interface{}, msg string, args ...interface{}) bool
```

The `.Done()` function lets the WaitGroup know it has completed. (Blocked until the WaitGroup counter goes back to 0). The WaitGroup at the end of the main test function is waiting for the **Done** message.

## Completion

The last part is the Test function runs a `assert.NoError(t, err)` check.

## Summary
The advantage of running unit tests allow us to mock situations to test if the code meets the requirements.
By running the unit tests before Pulumi Up is run allows us to *fail fast*.
