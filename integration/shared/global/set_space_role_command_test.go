package global

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"

	"code.cloudfoundry.org/cli/integration/helpers"
)

var _ = Describe("set-space-role command", func() {
	When("the set_roles_by_username flag is disabled", func() {
		BeforeEach(func() {
			helpers.LoginCF()
			helpers.DisableFeatureFlag("set_roles_by_username")
		})

		AfterEach(func() {
			helpers.EnableFeatureFlag("set_roles_by_username")
		})

		When("the user does not exist", func() {
			It("prints the error from UAA and exits 1", func() {
				session := helpers.CF("set-space-role", "not-exists", "some-org", "some-space", "SpaceDeveloper")
				Eventually(session).Should(Say("FAILED"))
				Eventually(session).Should(Say("User not-exists not found"))
				Eventually(session).Should(Exit(1))
			})
		})
	})
})
