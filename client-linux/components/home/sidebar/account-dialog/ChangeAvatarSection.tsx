import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Check, Loader2 } from "lucide-react";
import * as z from "zod";
import { Form } from "@/components/ui/form";
import { useSession } from "next-auth/react";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import AvatarInput from "@/components/auth/AvatarInput";
import useFetch from "@/hooks/useFetch";

const formSchema = z.object({
  avatar: z.any().refine((val) => {
    return (
      typeof FileList !== "undefined" &&
      val instanceof FileList &&
      val.length > 0
    );
  }),
});

const ChangeAvatarSection = () => {
  const { data, update } = useSession();
  const { isLoading, exec } = useFetch({
    path: "/user/avatar",
    method: "PUT",
    auth: true,
  });
  const form = useForm({
    resolver: zodResolver(formSchema),
    defaultValues: {
      avatar: undefined,
    },
  });

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    try {
      const formData = new FormData();
      formData.append("file", values.avatar[0]);

      const { ok, responseData, error } = await exec(formData);
      if (error) throw error;
      if (ok) {
        toast.success(responseData.message);
        await update({ avatar: responseData.avatar });
        form.reset();
      } else {
        toast.error(responseData.message);
      }
    } catch (err) {
      console.log(err);
      toast.error("Coudln't change your avatar, please try again later.");
    }
  };

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="relative w-fit mx-auto hover:opacity-60"
      >
        <AvatarInput
          control={form.control}
          name="avatar"
          current={`${process.env.AWS}${data?.user.avatar}`}
        />
        <div className="absolute rounded-full top-0 right-0 flex">
          {isLoading ? (
            <Button type="button" disabled className="rounded-full" size="sm">
              <Loader2 className="animate-spin" />
            </Button>
          ) : (
            form.formState.isValid && (
              <Button
                type="submit"
                className="rounded-full bg-green-600 hover:bg-green-600"
                size="sm"
              >
                <Check />
              </Button>
            )
          )}
        </div>
      </form>
    </Form>
  );
};

export default ChangeAvatarSection;
