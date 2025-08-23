import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { DropdownMenuItem } from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { Mail, User } from "lucide-react";
import { useSession } from "next-auth/react";
import ChangeNameSection from "./ChangeNameSection";
import { Label } from "@/components/ui/label";
import ChangePasswordSection from "./ChangePasswordSection";
import ChangeAvatarSection from "./ChangeAvatarSection";

const AccountDialog = () => {
  const { data } = useSession();

  return (
    <Dialog>
      <DialogTrigger asChild>
        <DropdownMenuItem onSelect={(e) => e.preventDefault()}>
          <User className="mr-2 h-4 w-4" />
          <span>Account</span>
        </DropdownMenuItem>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle className="text-center">Settings</DialogTitle>
        </DialogHeader>

        <ChangeAvatarSection />

        <ChangeNameSection />

        <div>
          <Label>Email</Label>
          <div className="relative">
            <Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              type="email"
              className="pl-10"
              disabled={true}
              value={data?.user.email}
            />
          </div>
        </div>

        <ChangePasswordSection />

        <DialogFooter>
          <DialogClose asChild>
            <Button variant="outline">Cancel</Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default AccountDialog;
