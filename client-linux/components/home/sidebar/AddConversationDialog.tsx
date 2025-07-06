"use client";

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
import { Input } from "@/components/ui/input";
import { Ban, Plus, Search } from "lucide-react";
import { useRef, useState } from "react";
import AddConversationResults from "./AddConversationResults";
import { toast } from "sonner";
import useFetch from "@/hooks/useFetch";

type CloseBtnType = HTMLButtonElement | null;

export default function AddConversationDialog() {
  const closeBtnRef = useRef<CloseBtnType>(null);
  const [term, setTerm] = useState<string>("");
  const [currentSearch, setCurrentSearch] = useState<string>("");
  const { isLoading, exec } = useFetch({
    path: `/users?search=${term}`,
    auth: true,
  });
  const [result, setResult] = useState<Participant[]>([]);

  async function onSubmitHandler(e: React.FormEvent) {
    e.preventDefault();
    if (!term || term == currentSearch) {
      return;
    }
    try {
      const { ok, responseData, error } = await exec();
      if (error) throw error;
      else if (!ok) {
        toast.error(responseData.message);
        return;
      }

      setCurrentSearch(term);
      setResult(responseData.users || []);
    } catch (error) {
      toast.error("Something wen wrong, please try again later.");
      console.log(error);
    }
  }

  function onClose() {
    closeBtnRef.current?.click();
  }

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button>
          <Plus />
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <form onSubmit={onSubmitHandler}>
          <DialogHeader className="py-5">
            <DialogTitle className="text-center">
              Create a new conversation
            </DialogTitle>
          </DialogHeader>
          <div className="grid gap-4 mb-10">
            <div className="flex gap-2">
              <Input
                id="name-1"
                name="name"
                placeholder="Search..."
                value={term}
                onChange={(e) => setTerm(e.target.value)}
              />
              <Button disabled={!term || isLoading} type="submit">
                <Search />
              </Button>
            </div>
          </div>
          {!currentSearch && (
            <DialogFooter className="flex gap-3 text-xl font-bold justify-center">
              Type something to search <Search />
            </DialogFooter>
          )}
          {currentSearch && result.length == 0 && (
            <DialogFooter className="flex gap-3 text-xl font-bold justify-center">
              No results found. <Ban />
            </DialogFooter>
          )}
          <div className="flex flex-col gap-4">
            {result.length > 0 &&
              result.map((user) => (
                <AddConversationResults
                  key={user._id}
                  user={user}
                  onClose={onClose}
                ></AddConversationResults>
              ))}
          </div>
        </form>
        <DialogClose ref={closeBtnRef} />
      </DialogContent>
    </Dialog>
  );
}
