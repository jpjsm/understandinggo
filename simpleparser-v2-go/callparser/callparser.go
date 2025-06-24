package main

import (
	"fmt"

	"example.com/booleanparser"
)

type ExpressionTest struct {
	universe [][]string
	context  []string
	cases    []struct {
		expression      string
		expected_result bool
		throw_error     bool
	}
}

func main() {
	expressiontst := ExpressionTest{
		universe: [][]string{
			{"Read", "4246b7a7-1e49-40dd-8fa6-7aebdd70f34d"},
			{"Update", "44379cdf-2521-42f9-904e-c31d7244ed6c"},
			{"Insert", "d31aeb5b-e357-4a50-9a0f-3dda18b632ff"},
			{"Delete", "aa1ee703-e889-4b0d-8fa3-a39118a3443e"},
			{"Create", "d341c9da-1f00-414b-9d02-d109f01d4f70"},
			{"Alter", "b0e88dc8-a852-4e32-b4f7-42da2bd170fe"},
			{"Execute", "a5b2a69b-d7d5-46bf-bce9-d1cdaca88f54"},
			{"Take_Ownership", "6c639a12-53fd-4575-abfb-0bd61913c2af"},
			{"Impersonate", "eb4caa4d-7931-4e9e-b223-59b60b827461"},
		},
		context: []string{
			"44379cdf-2521-42f9-904e-c31d7244ed6c", // Update
			"aa1ee703-e889-4b0d-8fa3-a39118a3443e", // Delete
			"b0e88dc8-a852-4e32-b4f7-42da2bd170fe", // Alter
			"6c639a12-53fd-4575-abfb-0bd61913c2af", // Take_Ownership
			"00000000-0000-0000-0000-000000000000", // Not in universe
		},
		cases: []struct {
			expression      string
			expected_result bool
			throw_error     bool
		}{
			{expression: "Read, Update, Insert, Delete, Create, Alter, Execute, Take_Ownership, Impersonate", expected_result: true},
			{expression: "Read | Update | Insert | Delete | Create | Alter | Execute | Take_Ownership | Impersonate", expected_result: true},
			{expression: "Take_Ownership", expected_result: true},
			{expression: "!Take_Ownership", expected_result: false},
			{expression: "00000000-0000-0000-0000-000000000000", expected_result: true},
			{expression: "Update&Delete&Alter", expected_result: true},
			{expression: "!(Update&Delete&Alter)", expected_result: false},
			{expression: "(Update | Insert) & !Execute)", expected_result: true},
			{expression: "Update+Delete*Alter", throw_error: true},
			{expression: "!(mañana * (pingüino,árbol,garçon))", throw_error: true},
		},
	}

	for _, _case := range expressiontst.cases {
		observedvalue, evaluationerror := booleanparser.EvaluateBooleanExpression(
			_case.expression,
			expressiontst.context,
			expressiontst.universe)
		observederror := evaluationerror != nil
		expectederror := _case.throw_error
		if expectederror != observederror {
			fmt.Printf("[ERROR # '%s'] Expected error '%v' != Observed error '%v' => Expected value '%v', Observed value '%v', evaluation-error '%v'\n",
				_case.expression, expectederror, observederror, _case.expected_result, observedvalue, evaluationerror)
		} else {
			if !expectederror {
				if observedvalue != _case.expected_result {
					fmt.Printf("[TEST FAILED    # '%s'] Expected value '%v' != Observed value '%v'\n", _case.expression, _case.expected_result, observedvalue)
				} else {
					fmt.Printf("[TEST SUCCEEDED # '%s'] Expected value '%v' == Observed value '%v'\n", _case.expression, _case.expected_result, observedvalue)
				}
			} else {
				fmt.Printf("[TEST SUCCEEDED # '%s'] Exception received '%s'\n", _case.expression, evaluationerror)
			}
		}
	}
}
