'use client'
import {useState} from "react";
import {Button} from "@/app/components/ui/button";
import {Input} from "@/app/components/ui/input";
import {Label} from "@/app/components/ui/label";
import MessageComponent from "@/app/components/message";
import {Message} from "@/api/api";
import Spinner from "@/app/components/spinner";
import {InitPayment} from "@/api/payment";
import {PlusCircle} from "lucide-react";
import SimpleModal from "@/app/components/simple-modal";

export default function TopUpModal() {
    const [amount, setAmount] = useState<string>("1000");
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState<Message | null>(null);
    const [isOpen, setIsOpen] = useState(false);

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
        <>
            <Button
                onClick={() => setIsOpen(true)}
                variant="outline"
                size="sm"
                className="ml-2 gap-2 border-[#DB935B] text-[#DB935B] hover:bg-[#DB935B] hover:text-black z-10 pointer-events-auto relative"
            >
                <PlusCircle className="w-4 h-4" /> Пополнить
            </Button>

            <SimpleModal
                isOpen={isOpen}
                onClose={() => setIsOpen(false)}
                title="Пополнение баланса"
            >
                <div className="space-y-4">
                    <div>
                        <Label htmlFor="topup-amount" className="text-sm">Сумма (₽)</Label>
                        <Input
                            id="topup-amount"
                            type="number"
                            value={amount}
                            min={100}
                            onChange={(e) => setAmount(e.target.value)}
                            className="bg-[#301EBD] border-none text-white mt-1"
                        />
                    </div>
                    <p className="text-xs text-gray-400">
                        Вы будете перенаправлены на страницу ЮKassa для оплаты банковской картой.
                    </p>
                    <MessageComponent message={message} />
                    <Button 
                        onClick={handleTopUp} 
                        disabled={loading} 
                        className="bg-[#B7FF00] text-black hover:bg-[#a6e600] font-bold w-full"
                    >
                        {loading ? <Spinner size={20} /> : `Перейти к оплате`}
                    </Button>
                </div>
            </SimpleModal>
        </>
    );
}