'use client'
import {useState} from "react";
import {Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger} from "@/app/components/ui/dialog";
import {Button} from "@/app/components/ui/button";
import {Input} from "@/app/components/ui/input";
import {Label} from "@/app/components/ui/label";
import MessageComponent from "@/app/components/message";
import {Message} from "@/api/api";
import Spinner from "@/app/components/spinner";
import {InitPayment} from "@/api/payment";
import {PlusCircle} from "lucide-react";

export default function TopUpModal() {
    const [amount, setAmount] = useState<string>("1000");
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState<Message | null>(null);

    const handleTopUp = async () => {
        setLoading(true);
        setMessage(null);

        // Return URL - текущая страница
        const returnUrl = window.location.href;

        const url = await InitPayment(parseFloat(amount), returnUrl, setMessage);

        if (url) {
            // Редирект на платежный шлюз
            window.location.href = url;
        } else {
            setLoading(false);
        }
    };

    return (
        <Dialog>
            <DialogTrigger asChild>
                <Button variant="outline" size="sm" className="ml-2 gap-2 border-light-green text-[#DB935B] hover:bg-light-green hover:text-black">
                    <PlusCircle className="w-4 h-4" /> Пополнить
                </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[400px] bg-[#656662] text-white border-gray-700">
                <DialogHeader>
                    <DialogTitle>Пополнение баланса</DialogTitle>
                </DialogHeader>

                <div className="grid gap-4 py-4">
                    <div className="grid gap-2">
                        <Label htmlFor="topup-amount">Сумма (₽)</Label>
                        <Input
                            id="topup-amount"
                            type="number"
                            value={amount}
                            min={100}
                            onChange={(e) => setAmount(e.target.value)}
                            className="bg-[#301EBD] border-none text-white text-lg"
                        />
                    </div>
                    <p className="text-xs text-gray-400">
                        Вы будете перенаправлены на страницу ЮKassa для оплаты банковской картой или СБП.
                    </p>
                </div>

                <MessageComponent message={message} />

                <Button onClick={handleTopUp} disabled={loading} className="bg-[#B7FF00] text-black hover:bg-[#a6e600] font-bold w-full">
                    {loading ? <Spinner size={24} /> : `Перейти к оплате`}
                </Button>
            </DialogContent>
        </Dialog>
    );
}