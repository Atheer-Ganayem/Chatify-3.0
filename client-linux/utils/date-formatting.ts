export function getDate(createdAt: Date) {
  const dayDeffirence = new Date().getDay() - createdAt.getDay();

  let date = createdAt.toLocaleDateString();

  if (dayDeffirence === 0) {
    date = "Today";
  } else if (dayDeffirence === 1) {
    date = "Yesterday";
  }

  return (
    date + " at " + createdAt.toLocaleTimeString("eng-US", { hour: "numeric", minute: "numeric" })
  );
}
