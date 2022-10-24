/*
 * Copyright (c) 2022 Red Hat, Inc.
 * SPDX-License-Identifier: GPL-2.0-or-later
 */
package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Psql Tests", func() {

	AfterEach(func() {
		defer GinkgoRecover()
	})

	When("notIn", func() {
		It("Should return false if a value is not in the slice", func() {
			data := []int{
				1,
				3,
				5,
			}

			res := notIn(data, 2)

			Expect(res).To(BeTrue())
		})

		It("Should return true if a value is in the slice", func() {
			data := []int{
				1,
				3,
				5,
			}

			res := notIn(data, 1)

			Expect(res).To(BeFalse())
		})

		It("Should return true if using an empty slice", func() {
			data := []int{}

			res := notIn(data, 1)

			Expect(res).To(BeTrue())
		})
	})

	When("removeDuplicate", func() {
		It("Should return an ordered but unmodified slice", func() {
			entries := []entry{
				{symId: 3},
				{symId: 1},
			}

			res := removeDuplicate(entries)

			Expect(len(res)).To(Equal(2))
			Expect(res).To(Equal([]entry{{symId: 1}, {symId: 3}}))
		})

		It("Should return a new slice with no duplicated values", func() {
			entries := []entry{
				{symId: 1},
				{symId: 3},
				{symId: 1},
			}

			res := removeDuplicate(entries)

			Expect(len(res)).To(Equal(2))
			Expect(res).To(Equal([]entry{{symId: 1}, {symId: 3}}))
		})

		It("Should return an empty slice", func() {
			entries := []entry{}

			res := removeDuplicate(entries)

			Expect(len(res)).To(Equal(0))
			Expect(res).To(BeNil())
		})
	})

	When("notExcluded", func() {
		excluded := []string{"sym1", "sym2"}

		It("Should return false if item is in excluded slice", func() {
			res := notExcluded("sym1", excluded)

			Expect(res).To(BeFalse())
		})

		It("Should return true if item is not in excluded slice", func() {
			res := notExcluded("sym3", excluded)

			Expect(res).To(BeTrue())
		})

		It("Should return true if slice is empty", func() {
			res := notExcluded("sym1", []string{})

			Expect(res).To(BeTrue())
		})
	})

	When("navigate", func() {
		// TODO: `psql.navigate` fn refactor needed
	})

	Describe("intargets", func() {
		targets := []string{"n1", "n2"}

		When("There is a target match", func() {
			It("Should return true if n1 is a match", func() {
				res := intargets(targets, "n1", "n3")

				Expect(res).To(BeTrue())
			})

			It("Should return true if n2 is a match", func() {
				res := intargets(targets, "n3", "n2")

				Expect(res).To(BeTrue())
			})

			It("Should return true if both are a match", func() {
				res := intargets(targets, "n1", "n2")

				Expect(res).To(BeTrue())
			})
		})

		When("There no target match", func() {
			It("Should return false if neither is a match", func() {
				res := intargets(targets, "n3", "n5")

				Expect(res).To(BeFalse())
			})

			It("Should return false if the target slice is empty", func() {
				res := intargets([]string{}, "n1", "n2")

				Expect(res).To(BeFalse())
			})
		})
	})

})
