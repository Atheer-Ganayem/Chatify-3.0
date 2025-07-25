import { Input } from "@/components/ui/input";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Check, Loader2, Pen, User, X } from "lucide-react";
import * as z from "zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { useState } from "react";
import { useSession } from "next-auth/react";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import useFetch from "@/hooks/useFetch";

const formSchema = z.object({
  name: z.string().min(3, "Name must be at least 3 characters"),
});

const ChangeNameSection = () => {
  const { data, update } = useSession();
  const { isLoading, exec } = useFetch({
    path: "/user/name",
    method: "PUT",
    auth: true,
  });
  const [isEditing, setIsEditing] = useState<boolean>(false);

  const form = useForm({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: data?.user.name || "",
    },
  });

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    try {
      const { ok, responseData, error } = await exec();
      if (error) throw error;

      if (ok) {
        toast.success(responseData.message);
        await update({ name: values.name });
        setIsEditing(false);
      } else {
        toast.error(responseData.message);
      }
    } catch (err) {
      console.log(err);
      toast.error("Coudln't change your username, please try again later.");
    }
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem className="relative">
              <FormLabel>Username</FormLabel>
              <FormControl>
                <div className="relative">
                  <User className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                  <Input
                    disabled={!isEditing || isLoading}
                    placeholder="John Doe"
                    className="pl-10"
                    {...field}
                  />
                </div>
              </FormControl>
              <FormMessage />
              <div className="absolute rounded-full top-0 right-0 flex">
                {isLoading ? (
                  <Button
                    type="button"
                    disabled
                    className="rounded-full"
                    size="sm"
                  >
                    <Loader2 className="animate-spin" />
                  </Button>
                ) : !isEditing ? (
                  <Button
                    type="button"
                    className="rounded-full"
                    size="sm"
                    onClick={(e) => {
                      e.preventDefault();
                      setIsEditing(true);
                    }}
                  >
                    <Pen />
                  </Button>
                ) : (
                  <>
                    <Button
                      type="submit"
                      className="rounded-full bg-green-600 hover:bg-green-600"
                      size="sm"
                    >
                      <Check />
                    </Button>
                    <Button
                      type="button"
                      className="rounded-full"
                      onClick={() => {
                        setIsEditing(false);
                      }}
                      size="sm"
                      variant="destructive"
                    >
                      <X />
                    </Button>
                  </>
                )}
              </div>
            </FormItem>
          )}
        />
      </form>
    </Form>
  );
};

export default ChangeNameSection;
