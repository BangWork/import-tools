export const  containsSubstring=(str, substr) =>{
  const regex = new RegExp(substr, "i");
  return regex.test(str);
}
