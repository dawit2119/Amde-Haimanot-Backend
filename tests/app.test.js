import { describe, expect, it } from "vitest";
import { TestDriver } from "testdriverai/vitest/hooks";

describe("Amde Haymanot App", () => {
  it("should navigate through the main public sections and registration", async (context) => {
    const testdriver = TestDriver(context);

    // 1. Visit homepage
    await testdriver.provision.chrome({ url: "https://amdehaymanot.com" });
    let assertResult = await testdriver.assert(
      "The Amde Haymanot Sunday School page is displayed with the hero section and Enroll Now button"
    );
    expect(assertResult).toBeTruthy();

    // 2. Navigate to About page
    await testdriver.find("About link in the top navigation bar").click();
    assertResult = await testdriver.assert(
      "The About page is displayed with heading 'About Amdehayimanot Sunday School' and Our Story section"
    );
    expect(assertResult).toBeTruthy();

    // 3. Navigate to Contact page
    await testdriver.find("Contact link in the top navigation bar").click();
    assertResult = await testdriver.assert(
      "The Contact Us page is displayed with contact information and Send Us a Message form"
    );
    expect(assertResult).toBeTruthy();

    // 4. Navigate to Registration page
    await testdriver.find("Register button in the top navigation bar").click();
    assertResult = await testdriver.assert(
      "The Create Account registration form is displayed with Full Name and Email Address fields"
    );
    expect(assertResult).toBeTruthy();
  });
});
