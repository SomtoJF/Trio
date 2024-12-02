# Trio

Trio is a multi-agent chat app. Imagine if you had two LLMs in a groupchat. Thats Trio in a nutshell.

## Important information

- Trio is currently running on a free tier of Vercel. This means that the app sleeps after 5 minutes of inactivity. This is why it might take a while to load when you first visit the site.

- Trio is being migrated into new repos for the frontend and backend. I don't enjoy developing the frontend and backend in a monorepo and nx isn't exactly the best for developer experience. So the frontend and backend will be in separate repos.

- The new repo is [here](https://github.com/somtojf/trio-client).

- The new backend repo is [here](https://github.com/somtojf/trio-server).

## Why Trio?

- One cool thing is to see if they can work together at all.

- Another would be to see the extent of conversations these AI models can have with each other. Could they be capable of innovation if they work together?

## Updates

- 2024-10-29: I figured out something interesting a few weeks ago. Essentially how to make the LLMs more accurate. One problem with LLMs is that they tend to hallucinate. With the development of the new models, this problem has been somewhat mitigated. However, there is still some amount of hallucination and it increases as the complexity of the prompt/task increases. **What I figured out is basically this. If you feed in wrong information (specifically a hallucination) into the LLM and ask it to correct itself, it's pretty good at correcting itself.** After working on this on and off for a few weeks, I finally pushed a new feature called a `Reflection Chat`. It's a chat that is designed to enhance the accuracy of the LLM responses by placing two agents in a chat and having them iteratively provide responses (and far more importantly, correct and critique each other's responses) until they converge on the correct answer.
