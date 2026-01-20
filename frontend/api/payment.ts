import {api, DefaultErrorHandler, Message} from "@/api/api";

export async function InitPayment(amount: number, returnUrl: string, setMessage: (msg: Message) => void): Promise<string | null> {
    try {

        const userStr = localStorage.getItem("user");
        const user = userStr ? JSON.parse(userStr) : null;

        if (!user) throw new Error("User not found");

        const res = await api.post("/payment/pay/init", {
            entity_type: "user", //TODO someday change it
            entity_id: user.id,
            amount: amount,
            return_url: returnUrl
        });

        return res.data.confirmation_url;
    } catch (e: any) {
        DefaultErrorHandler(setMessage)(e);
        return null;
    }
}