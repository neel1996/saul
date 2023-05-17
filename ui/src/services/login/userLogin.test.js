import { post } from "@root/AxiosInstance";

import { login } from "./userLogin";

jest.mock("@root/AxiosInstance", () => ({
    post: jest.fn(),
}));

describe("userLogin service", () => {
    it("should invoke login api and store auth token", async () => {
        post.mockResolvedValue({
            authToken: "testAuthToken",
        });

        const data = await login({
            name: "test",
            email: "test@test.com",
        });

        expect(post).toHaveBeenCalledWith("/login", {
            name: "test",
            email: "test@test.com",
        });
        expect(data).toBeTruthy();
        expect(data.authToken).toBe("testAuthToken");
        expect(sessionStorage.getItem("authToken")).toBe("testAuthToken");
    });

    it("should throw error if login api fails", async () => {
        post.mockRejectedValue(new Error("testError"));

        await expect(
            login({
                name: "test",
                email: "test@test.com",
            })
        ).rejects.toThrow("testError");
    });

    it("should throw error if auth token is missing", async () => {
        post.mockResolvedValue({});

        await expect(
            login({
                name: "test",
                email: "test@test.com",
            })
        ).rejects.toThrow("Invalid credentials");
    });
});
