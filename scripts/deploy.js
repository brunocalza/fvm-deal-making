const hre = require("hardhat");

async function main() {
  const Contract = await hre.ethers.getContractFactory("DealClient");
  const contract = await Contract.deploy();

  console.log(`Deployed to ${contract.address}`);
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
